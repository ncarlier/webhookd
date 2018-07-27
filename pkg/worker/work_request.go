package worker

import (
	"sync"
	"sync/atomic"
)

var workID uint64

// WorkRequest is a request of work for a worker
type WorkRequest struct {
	ID          uint64
	Name        string
	Script      string
	Payload     string
	Args        []string
	MessageChan chan []byte
	Timeout     int
	Terminated  bool
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
		Terminated:  false,
	}
}

// Terminate set work request as terminated
func (wr *WorkRequest) Terminate() {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()
	if !wr.Terminated {
		wr.Terminated = true
	}
}

// IsTerminated ask if the work request is terminated
func (wr *WorkRequest) IsTerminated() bool {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()
	return wr.Terminated
}
