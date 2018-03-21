package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ncarlier/webhookd/pkg/hook"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/tools"
	"github.com/ncarlier/webhookd/pkg/worker"
)

var (
	defaultTimeout = atoiFallback(os.Getenv("APP_HOOK_TIMEOUT"), 10)
)

func atoiFallback(str string, fallback int) int {
	value, err := strconv.Atoi(str)
	if err != nil || value < 0 {
		return fallback
	}
	return value
}

// Index is the main handler of the API.
func Index() http.Handler {
	return http.HandlerFunc(webhookHandler)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
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
		logger.Error.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	params := tools.QueryParamsToShellVars(r.URL.Query())
	logger.Debug.Printf("Calling hook script \"%s\" with params %s...\n", script, params)
	params = append(params, tools.HTTPHeadersToShellVars(r.Header)...)

	// Create work
	timeout := atoiFallback(r.Header.Get("X-Hook-Timeout"), defaultTimeout)
	work := worker.NewWorkRequest(p, script, string(body), params, timeout)

	// Put work in queue
	worker.WorkQueue <- *work

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	logger.Debug.Println("Work request queued:", script)
	// fmt.Fprintf(w, "data: Running \"%s\" ...\n\n", work.Name)

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
