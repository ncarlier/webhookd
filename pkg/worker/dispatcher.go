package worker

import (
	"github.com/ncarlier/webhookd/pkg/logger"
)

// WorkerQueue is the global queue of Workers
var WorkerQueue chan chan Work

// WorkQueue is the global queue of work to dispatch
var WorkQueue = make(chan Work, 100)

// StartDispatcher is charged to start n workers.
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan Work, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		logger.Debug.Printf("starting worker #%d ...\n", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue

					logger.Debug.Printf("dispatching hook request: %s#%d", work.Name(), work.ID())
					worker <- work
				}()
			}
		}
	}()
}
