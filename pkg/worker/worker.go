package worker

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/webhookd/pkg/metric"

	"github.com/ncarlier/webhookd/pkg/notification"
)

// NewWorker creates, and returns a new Worker object.
func NewWorker(id int, workerQueue chan chan Work) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan Work),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}

	return worker
}

// Worker is a go routine in charge of executing a work.
type Worker struct {
	ID          int
	Work        chan Work
	WorkerQueue chan chan Work
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
				slog.Debug("hook execution request received", "worker", w.ID, "hook", work.Name(), "id", work.ID())
				metric.Requests.Add(1)
				err := work.Run()
				if err != nil {
					metric.RequestsFailed.Add(1)
					work.SendMessage(fmt.Sprintf("error: %s", err.Error()))
				}
				// Send notification
				go notification.Notify(work)

				work.Close()
			case <-w.QuitChan:
				slog.Debug("stopping worker...", "worker", w.ID)
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
