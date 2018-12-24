package model

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/tools"
)

var workID uint64

// WorkStatus is the status of a worload
type WorkStatus int

const (
	// Idle means that the work is not yet started
	Idle WorkStatus = iota
	// Running means that the work is running
	Running
	// Success means that the work over
	Success
	// Error means that the work is over but in error
	Error
)

// WorkRequest is a request of work for a worker
type WorkRequest struct {
	ID          uint64
	Name        string
	Script      string
	Payload     string
	Args        []string
	MessageChan chan []byte
	Timeout     int
	Status      WorkStatus
	LogFilename string
	Err         error
	mutex       sync.Mutex
}

// NewWorkRequest creats new work request
func NewWorkRequest(name, script, payload string, args []string, timeout int) *WorkRequest {
	w := &WorkRequest{
		ID:          atomic.AddUint64(&workID, 1),
		Name:        name,
		Script:      script,
		Payload:     payload,
		Args:        args,
		Timeout:     timeout,
		MessageChan: make(chan []byte),
		Status:      Idle,
	}
	w.LogFilename = path.Join(*config.Get().LogDir, fmt.Sprintf("%s_%d_%s.txt", tools.ToSnakeCase(w.Name), w.ID, time.Now().Format("20060102_1504")))
	return w
}

// Terminate set work request as terminated
func (wr *WorkRequest) Terminate(err error) error {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()
	if err != nil {
		wr.Status = Error
		wr.Err = err
		logger.Info.Printf("Work %s#%d done [ERROR]\n", wr.Name, wr.ID)
		return err
	}
	wr.Status = Success
	logger.Info.Printf("Work %s#%d done [SUCCESS]\n", wr.Name, wr.ID)
	return nil
}

// IsTerminated ask if the work request is terminated
func (wr *WorkRequest) IsTerminated() bool {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()
	return wr.Status == Success || wr.Status == Error
}

// GetLogContent returns work logs filtered with the prefix
func (wr *WorkRequest) GetLogContent(prefixFilter string) string {
	file, err := os.Open(wr.LogFilename)
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
