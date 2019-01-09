package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/notification"
	"github.com/ncarlier/webhookd/pkg/worker"
)

func main() {
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	conf := config.Get()

	level := "info"
	if *conf.Debug {
		level = "debug"
	}
	logger.Init(level)

	logger.Debug.Println("Starting webhookd server...")

	server := &http.Server{
		Addr:     *conf.ListenAddr,
		Handler:  api.NewRouter(config.Get()),
		ErrorLog: logger.Error,
	}

	// Configure notification
	if err := notification.Init(*conf.NotificationURI); err != nil {
		logger.Error.Fatalf("Unable to create notification channel: %v\n", err)
	}

	// Start the dispatcher.
	logger.Debug.Printf("Starting the dispatcher (%d workers)...\n", *conf.NbWorkers)
	worker.StartDispatcher(*conf.NbWorkers)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Debug.Println("Server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Error.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Info.Println("Server is ready to handle requests at", *conf.ListenAddr)
	api.Start()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error.Fatalf("Could not listen on %s: %v\n", *conf.ListenAddr, err)
	}

	<-done
	logger.Debug.Println("Server stopped")
}
