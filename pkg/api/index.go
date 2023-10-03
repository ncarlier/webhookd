package api

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/hook"
	"github.com/ncarlier/webhookd/pkg/worker"
)

var (
	defaultTimeout int
	scriptDir      string
	outputDir      string
)

func atoiFallback(str string, fallback int) int {
	if value, err := strconv.Atoi(str); err == nil && value > 0 {
		return value
	}
	return fallback
}

// index is the main handler of the API.
func index(conf *config.Config) http.Handler {
	defaultTimeout = conf.HookTimeout
	scriptDir = conf.ScriptDir
	outputDir = conf.HookLogDir
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
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Get hook location
	hookName := strings.TrimPrefix(r.URL.Path, "/")
	if hookName == "" {
		infoHandler(w, r)
		return
	}
	_, err := hook.ResolveScript(scriptDir, hookName)
	if err != nil {
		slog.Error("hooke not found", "err", err.Error())
		http.Error(w, "hook not found", http.StatusNotFound)
		return
	}

	if err = r.ParseForm(); err != nil {
		slog.Error("error reading from-data", "err", err)
		http.Error(w, "unable to parse request form", http.StatusBadRequest)
		return
	}

	// parse body
	var body []byte
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediatype, _, _ := mime.ParseMediaType(ct)
		if strings.HasPrefix(mediatype, "text/") || mediatype == "application/json" {
			body, err = io.ReadAll(r.Body)
			if err != nil {
				slog.Error("error reading body", "err", err)
				http.Error(w, "unable to read request body", http.StatusBadRequest)
				return
			}
		}
	}

	params := HTTPParamsToShellVars(r.Form)
	params = append(params, HTTPParamsToShellVars(r.Header)...)

	// Create work
	timeout := atoiFallback(r.Header.Get("X-Hook-Timeout"), defaultTimeout)
	job, err := hook.NewHookJob(&hook.Request{
		Name:      hookName,
		Method:    r.Method,
		Payload:   string(body),
		Args:      params,
		Timeout:   timeout,
		BaseDir:   scriptDir,
		OutputDir: outputDir,
	})
	if err != nil {
		slog.Error("error creating hook job", "err", err)
		http.Error(w, "unable to create hook job", http.StatusInternalServerError)
		return
	}

	// Put work in queue
	worker.WorkQueue <- job

	// Use content negotiation to enable Server-Sent Events
	useSSE := r.Method == "GET" && r.Header.Get("Accept") == "text/event-stream"
	if useSSE {
		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
	} else {
		// Send chunked response
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Hook-ID", strconv.FormatUint(job.ID(), 10))

	for {
		msg, open := <-job.MessageChan
		if !open {
			break
		}
		if useSSE {
			fmt.Fprintf(w, "data: %s\n\n", msg) // Send SSE response
		} else {
			fmt.Fprintf(w, "%s\n", msg) // Send chunked response
		}
		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()
	}
}

func getWebhookLog(w http.ResponseWriter, r *http.Request) {
	// Get hook ID
	id := path.Base(r.URL.Path)

	// Get script location
	hookName := path.Dir(strings.TrimPrefix(r.URL.Path, "/"))
	_, err := hook.ResolveScript(scriptDir, hookName)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Retrieve log file
	logFile, err := hook.Logs(id, hookName, outputDir)
	if err != nil {
		slog.Error(err.Error())
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
