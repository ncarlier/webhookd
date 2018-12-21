package worker

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// ChanWriter is a simple writer to a channel of byte.
type ChanWriter struct {
	ByteChan chan []byte
}

func (c *ChanWriter) Write(p []byte) (int, error) {
	c.ByteChan <- p
	return len(p), nil
}

func run(work *WorkRequest) error {
	work.Status = Running
	logger.Info.Printf("Work %s#%d started...\n", work.Name, work.ID)
	logger.Debug.Printf("Work %s#%d script: %s\n", work.Name, work.ID, work.Script)
	logger.Debug.Printf("Work %s#%d parameter: %v\n", work.Name, work.ID, work.Args)

	binary, err := exec.LookPath(work.Script)
	if err != nil {
		return work.Terminate(err)
	}

	// Exec script with args...
	cmd := exec.Command(binary, work.Payload)
	// with env variables...
	cmd.Env = append(os.Environ(), work.Args...)
	// using a process group...
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Open the log file for writing
	logFile, err := createLogFile(work)
	if err != nil {
		return work.Terminate(err)
	}
	defer logFile.Close()
	logger.Debug.Printf("Work %s#%d output to file: %s\n", work.Name, work.ID, logFile.Name())

	wLogFile := bufio.NewWriter(logFile)
	defer wLogFile.Flush()

	// Combine cmd stdout and stderr
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return work.Terminate(err)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return work.Terminate(err)
	}
	cmdReader := io.MultiReader(outReader, errReader)

	// Start the script...
	err = cmd.Start()
	if err != nil {
		return work.Terminate(err)
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
			if !work.IsTerminated() {
				work.MessageChan <- []byte(line)
			} else {
				logger.Error.Printf("Work %s#%d is over. Unable to write more data into the channel: %s\n", work.Name, work.ID, line)
				break
			}
			// writing to outfile
			if _, err := wLogFile.WriteString(line + "\n"); err != nil {
				logger.Error.Println("Error while writing into the log file:", logFile.Name(), err)
				break
			}
		}
		if err := scanner.Err(); err != nil {
			logger.Error.Printf("Work %s#%d unable to read script stdout: %v\n", work.Name, work.ID, err)
		}
		wg.Done()
	}(cmdReader)

	// Start timeout timer
	timer := time.AfterFunc(time.Duration(work.Timeout)*time.Second, func() {
		logger.Warning.Printf("Work %s#%d has timed out (%ds). Killing process #%d...\n", work.Name, work.ID, work.Timeout, cmd.Process.Pid)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	// Wait for command output completion
	wg.Wait()

	// Wait for command completion
	err = cmd.Wait()

	// Stop timeout timer
	timer.Stop()

	// Mark work as terminated
	return work.Terminate(err)
}
