package main

import (
	"flag"
	"github.com/ncarlier/webhookd/api"
	"github.com/ncarlier/webhookd/worker"
	"log"
	"net/http"
)

var (
	LAddr    = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
	NWorkers = flag.Int("n", 2, "The number of workers to start")
)

func main() {
	flag.Parse()
	
	log.Println("Starting webhookd server...")

	// Start the dispatcher.
	log.Println("Starting the dispatcher")
	worker.StartDispatcher(*NWorkers)

	log.Println("Starting the http server")
	log.Fatal(http.ListenAndServe(*LAddr, api.Handlers()))
}
