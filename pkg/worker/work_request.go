package worker

import (
	"sync"
	"sync/atomic"

	"github.com/ncarlier/webhookd/pkg/logger"
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
	mutex       sync.Mutex
}

// NewWorkRequest creats new work request
func NewWorkRequest(name, script, payload string, args []string, timeout int) *WorkRequest {
	return &WorkRequest{
		ID:          atomic.AddUint64(&workID, 1),
		Name:        name,
		Script:      script,
		Payload:     payload,
		Args:        args,
		Timeout:     timeout,
		MessageChan: make(chan []byte),
		Status:      Idle,
	}
}

// Terminate set work request as terminated
func (wr *WorkRequest) Terminate(err error) error {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()
	if err != nil {
		wr.Status = Error
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
