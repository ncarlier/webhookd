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
	"github.com/ncarlier/webhookd/pkg/helper"
	"github.com/ncarlier/webhookd/pkg/hook"
	"github.com/ncarlier/webhookd/pkg/worker"
)

var (
	defaultTimeout int
	defaultExt     string
	scriptDir      string
	outputDir      string
)

var supportedContentTypes = []string{"text/plain", "text/event-stream", "application/json", "text/*"}

func atoiFallback(str string, fallback int) int {
	if value, err := strconv.Atoi(str); err == nil && value > 0 {
		return value
	}
	return fallback
}

// index is the main handler of the API.
func index(conf *config.Config) http.Handler {
	defaultTimeout = conf.Hook.Timeout
	defaultExt = conf.Hook.DefaultExt
	scriptDir = conf.Hook.ScriptsDir
	outputDir = conf.Hook.LogDir
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
	script, err := hook.ResolveScript(scriptDir, hookName, defaultExt)
	if err != nil {
		msg := "hook not found"
		slog.Error(msg, "err", err.Error())
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	if err = r.ParseForm(); err != nil {
		msg := "unable to parse form-data"
		slog.Error(msg, "err", err)
		http.Error(w, msg, http.StatusBadRequest)
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
				msg := "unable to read request body"
				slog.Error(msg, "err", err)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
	}

	params := HTTPParamsToShellVars(r.Form)
	params = append(params, HTTPParamsToShellVars(r.Header)...)

	// Create hook job
	timeout := atoiFallback(r.Header.Get("X-Hook-Timeout"), defaultTimeout)
	job, err := hook.NewHookJob(&hook.Request{
		Name:      hookName,
		Script:    script,
		Method:    r.Method,
		Payload:   string(body),
		Args:      params,
		Timeout:   timeout,
		OutputDir: outputDir,
	})
	if err != nil {
		msg := "unable to create hook execution job"
		slog.Error(msg, "err", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Put work in queue
	worker.WorkQueue <- job

	// Use content negotiation
	ct = helper.NegotiateContentType(r, supportedContentTypes, "text/plain")

	// set respons headers
	w.Header().Set("Content-Type", ct+"; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Hook-ID", strconv.FormatUint(job.ID(), 10))

	for {
		msg, open := <-job.MessageChan
		if !open {
			break
		}
		if ct == "text/event-stream" {
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
	_, err := hook.ResolveScript(scriptDir, hookName, defaultExt)
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
		http.Error(w, "hook execution log not found", http.StatusNotFound)
		return
	}
	defer logFile.Close()

	w.Header().Set("Content-Type", "text/plain")

	io.Copy(w, logFile)
}
