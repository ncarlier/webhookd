package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
	"github.com/ncarlier/webhookd/pkg/tools"
	"github.com/ncarlier/webhookd/pkg/worker"
)

var (
	defaultTimeout int
	scriptDir      string
)

func atoiFallback(str string, fallback int) int {
	if value, err := strconv.Atoi(str); err == nil && value > 0 {
		return value
	}
	return fallback
}

// index is the main handler of the API.
func index(conf *config.Config) http.Handler {
	defaultTimeout = *conf.Timeout
	scriptDir = *conf.ScriptDir
	return http.HandlerFunc(webhookHandler)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if _, err := strconv.Atoi(filepath.Base(r.URL.Path)); err == nil {
			getWebhookLog(w, r)
			return
		}
	}
	triggerWebhook(w, r)
}

func triggerWebhook(w http.ResponseWriter, r *http.Request) {
	// Check that streaming is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	// Get script location
	p := strings.TrimPrefix(r.URL.Path, "/")
	script, err := tools.ResolveScript(scriptDir, p)
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
	params = append(params, tools.HTTPHeadersToShellVars(r.Header)...)

	// logger.Debug.Printf("API REQUEST: \"%s\" with params %s...\n", p, params)

	// Create work
	timeout := atoiFallback(r.Header.Get("X-Hook-Timeout"), defaultTimeout)
	work := model.NewWorkRequest(p, script, string(body), params, timeout)

	// Put work in queue
	worker.WorkQueue <- *work

	if r.Method == "GET" {
		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
	} else {
		// Send chunked response
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Hook-ID", strconv.FormatUint(work.ID, 10))

	for {
		msg, open := <-work.MessageChan

		if !open {
			break
		}

		if r.Method == "GET" {
			fmt.Fprintf(w, "data: %s\n\n", msg) // Send SSE response
		} else {
			fmt.Fprintf(w, "%s\n", msg) // Send chunked response
		}

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

func getWebhookLog(w http.ResponseWriter, r *http.Request) {
	// Get hook ID
	id := path.Base(r.URL.Path)

	// Get script location
	name := path.Dir(strings.TrimPrefix(r.URL.Path, "/"))
	_, err := tools.ResolveScript(scriptDir, name)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Retrieve log file
	logFile, err := worker.RetrieveLogFile(id, name)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if logFile == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	defer logFile.Close()

	w.Header().Set("Content-Type", "text/plain")

	io.Copy(w, logFile)
}
