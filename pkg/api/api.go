package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ncarlier/webhookd/pkg/hook"
	"github.com/ncarlier/webhookd/pkg/tools"
	"github.com/ncarlier/webhookd/pkg/worker"
)

// WebhookHandler is the main handler of the API.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get script location
	p := strings.TrimPrefix(r.URL.Path, "/")
	script, err := hook.ResolveScript(p)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	params := tools.QueryParamsToShellVars(r.URL.Query())
	log.Printf("Calling hook script \"%s\" with params %s...\n", script, params)

	// Create work
	work := new(worker.WorkRequest)
	work.Name = p
	work.Script = script
	work.Payload = string(body)
	work.Args = params
	work.MessageChan = make(chan []byte)

	// Put work in queue
	worker.WorkQueue <- *work

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("Work request queued:", script)
	fmt.Fprintf(w, "data: Hook work request \"%s\" queued...\n\n", work.Name)

	for {
		msg, open := <-work.MessageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}
