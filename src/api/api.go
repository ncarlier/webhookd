package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ncarlier/webhookd/hook"
	"github.com/ncarlier/webhookd/worker"
	"log"
	"net/http"
)

func createWebhookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hookname := params["hookname"]
	action := params["action"]

	// Get hook decoder
	record, err := hook.RecordFactory(hookname)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("Using hook %s with action %s.\n", hookname, action)

	// Decode request
	err = record.Decode(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create work
	work := new(worker.WorkRequest)
	work.Name = hookname
	work.Action = action
	fmt.Println("Extracted data: ", record.GetURL(), record.GetName())
	work.Args = []string{record.GetURL(), record.GetName()}

	//Put work in queue
	worker.WorkQueue <- *work
	fmt.Printf("Work request queued: %s/%s\n", hookname, action)

	fmt.Fprintf(w, "Action %s of hook %s queued.", action, hookname)
}

func Handlers() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/{hookname:[a-z]+}/{action:[a-z]+}", createWebhookHandler).Methods("POST")
	return r
}

