package api

import (
	"bytes"
	"container/ring"
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
	defaultMode    string
	scriptDir      string
	outputDir      string
)

const (
	DefaultBufferLength = 100
	MaxBufferLength     = 10000
	SSEContentType      = "text/event-stream"
)

var supportedContentTypes = []string{"text/plain", SSEContentType, "application/json", "text/*"}

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
	defaultMode = conf.Hook.DefaultMode
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
	// Manage content negotiation
	negociatedContentType := helper.NegotiateContentType(r, supportedContentTypes, "text/plain")

	// Extract streaming method
	mode := r.Header.Get("X-Hook-Mode")
	if mode != "buffered" && mode != "chunked" {
		mode = defaultMode
	}
	if negociatedContentType == SSEContentType {
		mode = "sse"
	}

	// Check that streaming is supported
	if _, ok := w.(http.Flusher); !ok && mode != "buffered" {
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
		switch {
		case mediatype == "application/json", strings.HasPrefix(mediatype, "text/"):
			body, err = io.ReadAll(r.Body)
			if err != nil {
				msg := "unable to read request body"
				slog.Error(msg, "err", err)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		case mediatype == "multipart/form-data":
			if err := r.ParseMultipartForm(8 << 20); err != nil {
				msg := "unable to parse multipart/form-data"
				slog.Error(msg, "err", err)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		default:
			slog.Debug("unsuported media type", "media_type", mediatype)
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

	// Write hook ouput to the response regarding the asked method
	if mode != "buffered" {
		// Write hook response as Server Sent Event stream
		writeStreamedResponse(w, negociatedContentType, job, mode)
	} else {
		maxBufferLength := atoiFallback(r.Header.Get("X-Hook-MaxBufferedLines"), DefaultBufferLength)
		if maxBufferLength > MaxBufferLength {
			maxBufferLength = MaxBufferLength
		}
		// Write hook response after hook execution
		writeStandardResponse(w, negociatedContentType, job, maxBufferLength)
	}
}

func writeStreamedResponse(w http.ResponseWriter, negociatedContentType string, job *hook.Job, mode string) {
	writeHeaders(w, negociatedContentType, job.ID())
	for {
		msg, open := <-job.MessageChan
		if !open {
			break
		}

		if mode == "sse" {
			// Send SSE response
			prefix := "data: "
			if bytes.HasPrefix(msg, []byte("error:")) {
				prefix = ""
			}
			fmt.Fprintf(w, "%s%s\n", prefix, msg)
		} else {
			// Send chunked response
			w.Write(msg)
		}

		// Flush the data immediately instead of buffering it for later.
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

func writeStandardResponse(w http.ResponseWriter, negociatedContentType string, job *hook.Job, maxBufferLength int) {
	buffer := ring.New(maxBufferLength)
	overflow := false
	lines := 0

	// Consume messages into a ring buffer
	for {
		msg, open := <-job.MessageChan
		if !open {
			break
		}
		buffer.Value = msg
		buffer = buffer.Next()
		lines++
		if lines > maxBufferLength {
			overflow = true
		}
	}

	writeHeaders(w, negociatedContentType, job.ID())
	w.WriteHeader(getJobStatusCode(job))
	if overflow {
		w.Write([]byte("[output truncated]\n"))
	}
	// Write buffer to HTTP response
	buffer.Do(func(data interface{}) {
		if data != nil {
			w.Write(data.([]byte))
		}
	})
}

func getJobStatusCode(job *hook.Job) int {
	switch {
	case job.ExitCode() == 0:
		return http.StatusOK
	case job.ExitCode() >= 100:
		return job.ExitCode() + 300
	default:
		return http.StatusInternalServerError
	}
}

func writeHeaders(w http.ResponseWriter, contentType string, hookId uint64) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Hook-ID", strconv.FormatUint(hookId, 10))
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
	logFile, err := hook.GetLogFile(id, hookName, outputDir)
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
