package worker

import (
	"fmt"

	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/notification"
	"github.com/ncarlier/webhookd/pkg/tools"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Worker is a go routine in charge of executing a work.
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
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
				filename, err := run(&work)
				if err != nil {
					subject := fmt.Sprintf("Webhook %s#%d FAILED.", work.Name, work.ID)
					work.MessageChan <- []byte(fmt.Sprintf("error: %s", err.Error()))
					notify(subject, err.Error(), filename)
				} else {
					subject := fmt.Sprintf("Webhook %s#%d SUCCEEDED.", work.Name, workID)
					work.MessageChan <- []byte("done")
					notify(subject, "See attachment.", filename)
				}
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

func notify(subject string, text string, outfilename string) {
	var notifier, err = notification.NotifierFactory()
	if err != nil {
		logger.Debug.Println("Unable to get the notifier. Notification skipped:", err)
		return
	}
	if notifier == nil {
		logger.Debug.Println("Notification provider not found. Notification skipped.")
		return
	}

	var zipfile string
	if outfilename != "" {
		zipfile, err = tools.CompressFile(outfilename)
		if err != nil {
			fmt.Println(err)
			zipfile = outfilename
		}
	}

	notifier.Notify(subject, text, zipfile)
}
