package worker

import (
	"log/slog"
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
		slog.Debug("starting worker...", "worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			work := <-WorkQueue
			go func() {
				worker := <-WorkerQueue
				slog.Debug("dispatching hook request", "hook", work.Name(), "id", work.ID())
				worker <- work
			}()

		}
	}()
}
