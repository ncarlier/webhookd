package worker

import (
	"fmt"
	"github.com/ncarlier/webhookd/notification"
	"github.com/ncarlier/webhookd/tools"
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

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request %s/%s\n", w.ID, work.Name, work.Action)
				filename, err := RunScript(&work)
				if err != nil {
					subject := fmt.Sprintf("Webhook %s/%s FAILED.", work.Name, work.Action)
					Notify(subject, err.Error(), filename)
				} else {
					subject := fmt.Sprintf("Webhook %s/%s SUCCEEDED.", work.Name, work.Action)
					Notify(subject, "See attachment.", filename)
				}
			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func Notify(subject string, text string, outfilename string) {
	var notifier, err = notification.NotifierFactory()
	if err != nil {
		fmt.Println(err)
		return
	}
	if notifier == nil {
		fmt.Println("Notification provider not found.")
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
