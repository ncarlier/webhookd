package worker

import (
	"fmt"

	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
	"github.com/ncarlier/webhookd/pkg/notification"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan model.WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan model.WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Worker is a go routine in charge of executing a work.
type Worker struct {
	ID          int
	Work        chan model.WorkRequest
	WorkerQueue chan chan model.WorkRequest
	QuitChan    chan bool
}

// Start is the function to starts the worker by starting a goroutine.
// That is an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				logger.Debug.Printf("Worker #%d received work request: %s#%d\n", w.ID, work.Name, work.ID)
				err := run(&work)
				if err != nil {
					work.MessageChan <- []byte(fmt.Sprintf("error: %s", err.Error()))
				} else {
					work.MessageChan <- []byte("done")
				}
				// Send notification
				notification.Notify(&work)

				close(work.MessageChan)
			case <-w.QuitChan:
				logger.Debug.Printf("Stopping worker #%d...\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
