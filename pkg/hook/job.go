package hook

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ncarlier/webhookd/pkg/helper"
	"github.com/ncarlier/webhookd/pkg/logger"
)

var hookID uint64

// Job a hook job
type Job struct {
	id          uint64
	name        string
	script      string
	method      string
	payload     string
	args        []string
	MessageChan chan []byte
	timeout     int
	start       time.Time
	status      Status
	logFilename string
	err         error
	mutex       sync.Mutex
}

// NewHookJob creates new hook job
func NewHookJob(request *Request) (*Job, error) {
	script, err := ResolveScript(request.BaseDir, request.Name)
	if err != nil {
		return nil, err
	}
	job := &Job{
		id:          atomic.AddUint64(&hookID, 1),
		name:        request.Name,
		script:      script,
		method:      request.Method,
		payload:     request.Payload,
		args:        request.Args,
		timeout:     request.Timeout,
		MessageChan: make(chan []byte),
		status:      Idle,
	}
	job.logFilename = path.Join(request.OutputDir, fmt.Sprintf("%s_%d_%s.txt", helper.ToSnake(job.name), job.id, time.Now().Format("20060102_1504")))
	return job, nil
}

// ID returns job ID
func (job *Job) ID() uint64 {
	return job.id
}

// Name returns job name
func (job *Job) Name() string {
	return job.name
}

// Err returns job error
func (job *Job) Err() error {
	return job.err
}

// Meta returns job meta
func (job *Job) Meta() []string {
	return []string{
		"hook_id=" + strconv.FormatUint(job.id, 10),
		"hook_name=" + job.name,
		"hook_method=" + job.method,
	}
}

// Terminate set job as terminated
func (job *Job) Terminate(err error) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()
	job.status = Success
	if err != nil {
		job.status = Error
		job.err = err
		slog.Error(
			"hook executed",
			"hook", job.Name(),
			"id", job.ID(),
			"status", "error",
			"err", err,
			"took", time.Since(job.start).Microseconds(),
		)
		return err
	}
	slog.Info(
		"hook executed",
		"hook", job.Name(),
		"id", job.ID(),
		"status", "success",
		"took", time.Since(job.start).Microseconds(),
	)
	return nil
}

// IsTerminated ask if the job is terminated
func (job *Job) IsTerminated() bool {
	job.mutex.Lock()
	defer job.mutex.Unlock()
	return job.status == Success || job.status == Error
}

// Status get job status
func (job *Job) Status() Status {
	return job.status
}

// StatusLabel return job status as string
func (job *Job) StatusLabel() string {
	switch job.status {
	case Error:
		return "error"
	case Success:
		return "success"
	case Running:
		return "running"
	default:
		return "idle"
	}
}

// SendMessage send message to the message channel
func (job *Job) SendMessage(message string) {
	job.MessageChan <- []byte(message)
}

// Logs returns job logs filtered with the prefix
func (job *Job) Logs(prefixFilter string) string {
	file, err := os.Open(job.logFilename)
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	var result bytes.Buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefixFilter) {
			line = strings.TrimPrefix(line, prefixFilter)
			line = strings.TrimLeft(line, " ")
			result.WriteString(line + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err.Error()
	}
	return result.String()
}

// Close job message chan
func (job *Job) Close() {
	close(job.MessageChan)
}

// Run hook job
func (job *Job) Run() error {
	if job.status != Idle {
		return fmt.Errorf("unable to run job: status=%s", job.StatusLabel())
	}
	job.status = Running
	job.start = time.Now()
	slog.Info("executing hook...", "hook", job.name, "id", job.id)

	binary, err := exec.LookPath(job.script)
	if err != nil {
		return job.Terminate(err)
	}

	// Exec script with parameter...
	cmd := exec.Command(binary, job.payload)
	// with env variables and hook arguments...
	cmd.Env = append(os.Environ(), job.args...)
	// and hook meta...
	cmd.Env = append(cmd.Env, job.Meta()...)
	// using a process group...
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Open the log file for writing
	logFile, err := os.Create(job.logFilename)
	if err != nil {
		return job.Terminate(err)
	}
	defer logFile.Close()
	slog.Debug("hook details", "hook", job.name, "id", job.id, "script", job.script, "args", job.args, "output", logFile.Name())

	wLogFile := bufio.NewWriter(logFile)
	defer wLogFile.Flush()

	// Combine cmd stdout and stderr
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return job.Terminate(err)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return job.Terminate(err)
	}
	cmdReader := io.MultiReader(outReader, errReader)

	// Start the script...
	err = cmd.Start()
	if err != nil {
		return job.Terminate(err)
	}

	// Create wait group to wait for command output completion
	var wg sync.WaitGroup
	wg.Add(1)

	// Write script output to log file and the work message channel
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			// writing to the work channel
			if !job.IsTerminated() {
				job.MessageChan <- []byte(line)
			} else {
				slog.Error("hook execution done ; unable to write more data into the channel", "hook", job.name, "id", job.id, "line", line)
				break
			}
			// write to stdout if configured
			logger.LogIf(
				logger.HookOutputEnabled,
				slog.LevelInfo+1,
				line,
				"hook", job.name,
				"id", job.id,
			)
			// writing to outfile
			if _, err := wLogFile.WriteString(line + "\n"); err != nil {
				slog.Error("error while writing into the log file", "filename", logFile.Name(), "err", err)
				break
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Error("hook is unable to read script stdout", "hook", job.name, "id", job.id, "err", err)
		}
		wg.Done()
	}(cmdReader)

	// Start timeout timer
	timer := time.AfterFunc(time.Duration(job.timeout)*time.Second, func() {
		slog.Warn("hook has timed out: killing process...", "hook", job.name, "id", job.id, "timeout", job.timeout, "pid", cmd.Process.Pid)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	// Wait for command output completion
	wg.Wait()

	// Wait for command completion
	err = cmd.Wait()

	// Stop timeout timer
	timer.Stop()

	// Mark work as terminated
	return job.Terminate(err)
}
