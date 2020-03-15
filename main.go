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
	configflag "github.com/ncarlier/webhookd/pkg/config/flag"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/notification"
	"github.com/ncarlier/webhookd/pkg/server"
	"github.com/ncarlier/webhookd/pkg/version"
	"github.com/ncarlier/webhookd/pkg/worker"
)

func main() {
	conf := &config.Config{}
	configflag.Bind(conf, "WHD")

	flag.Parse()

	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	level := "info"
	if conf.Debug {
		level = "debug"
	}
	logger.Init(level)

	if conf.LogDir == "" {
		conf.LogDir = os.TempDir()
	}

	logger.Debug.Println("starting webhookd server...")

	srv := server.NewServer(conf)

	// Configure notification
	if err := notification.Init(conf.NotificationURI); err != nil {
		logger.Error.Fatalf("unable to create notification channel: %v\n", err)
	}

	// Start the dispatcher.
	logger.Debug.Printf("starting the dispatcher with %d workers...\n", conf.NbWorkers)
	worker.StartDispatcher(conf.NbWorkers)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Debug.Println("server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error.Fatalf("could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	addr := conf.ListenAddr
	if conf.TLSListenAddr != "" {
		addr = conf.TLSListenAddr
	}
	logger.Info.Println("server is ready to handle requests at", addr)
	api.Start()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error.Fatalf("could not listen on %s : %v\n", addr, err)
	}

	<-done
	logger.Debug.Println("server stopped")
}
