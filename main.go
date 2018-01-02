package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/worker"
)

var (
	lAddr     = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
	nbWorkers = flag.Int("n", 2, "The number of workers to start")
)

func main() {
	flag.Parse()

	log.Println("Starting webhookd server...")

	// Start the dispatcher.
	log.Printf("Starting the dispatcher (%d workers)...\n", *nbWorkers)
	worker.StartDispatcher(*nbWorkers)

	log.Printf("Starting the http server (%s)\n", *lAddr)
	http.HandleFunc("/", api.WebhookHandler)
	log.Fatal(http.ListenAndServe(*lAddr, nil))
}
