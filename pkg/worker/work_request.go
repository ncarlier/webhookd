package worker

import "sync/atomic"

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
	Closed      bool
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
		Closed:      false,
	}
}
