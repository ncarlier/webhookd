package worker

import (
	"github.com/ncarlier/webhookd/pkg/logger"
)

var WorkerQueue chan chan WorkRequest
var WorkQueue = make(chan WorkRequest, 100)

// StartDispatcher is charged to start n workers.
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		logger.Debug.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue

					logger.Debug.Printf("Dispatching work request: %s#%d", work.Name, work.ID)
					worker <- work
				}()
			}
		}
	}()
}
