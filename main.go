package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/auth"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/worker"
)

type key int

const (
	requestIDKey key = 0
)

var (
	healthy int32
)

func main() {
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	var authmethod auth.Method
	name := *config.Authentication
	if _, ok := auth.AvailableMethods[name]; ok {
		authmethod = auth.AvailableMethods[name]
		if err := authmethod.ParseParam(*config.AuthenticationParam); err != nil {
			fmt.Println("Authentication parameter is not valid:", err.Error())
			fmt.Println(authmethod.Usage())
			os.Exit(2)
		}
	} else {
		fmt.Println("Authentication name is not valid:", name)
		os.Exit(2)
	}

	level := "info"
	if *config.Debug {
		level = "debug"
	}
	logger.Init(level)

	logger.Debug.Println("Starting webhookd server...")
	logger.Debug.Println("Using Authentication:", name)
	authmethod.Init(*config.Debug)

	router := http.NewServeMux()
	router.Handle("/", api.Index(*config.Timeout, *config.ScriptDir))
	router.Handle("/healthz", healthz())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:     *config.ListenAddr,
		Handler:  authmethod.Middleware()(tracing(nextRequestID)(logging(logger.Debug)(router))),
		ErrorLog: logger.Error,
	}

	// Start the dispatcher.
	logger.Debug.Printf("Starting the dispatcher (%d workers)...\n", *config.NbWorkers)
	worker.StartDispatcher(*config.NbWorkers)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Debug.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Error.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Info.Println("Server is ready to handle requests at", *config.ListenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error.Fatalf("Could not listen on %s: %v\n", *config.ListenAddr, err)
	}

	<-done
	logger.Debug.Println("Server stopped")
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
