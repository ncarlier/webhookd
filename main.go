package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
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
	_ "github.com/ncarlier/webhookd/pkg/notification/all"
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

	if conf.HookLogDir == "" {
		conf.HookLogDir = os.TempDir()
	}

	if err := conf.Validate(); err != nil {
		log.Fatal("invalid configuration:", err)
	}

	logger.Configure(conf.LogFormat, conf.LogLevel)
	logger.HookOutputEnabled = conf.LogHookOutput
	logger.RequestOutputEnabled = conf.LogHTTPRequest

	slog.Debug("starting webhookd server...")

	srv := server.NewServer(conf)

	// Configure notification
	if err := notification.Init(conf.NotificationURI); err != nil {
		slog.Error("unable to create notification channel", "err", err)
	}

	// Start the dispatcher.
	slog.Debug("starting the dispatcher...", "workers", conf.NbWorkers)
	worker.StartDispatcher(conf.NbWorkers)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Debug("server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("could not gracefully shutdown the server", "err", err)
		}
		close(done)
	}()

	api.Start()
	slog.Info("server started", "addr", conf.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("unable to start the server", "addr", conf.ListenAddr, "err", err)
		os.Exit(1)
	}

	<-done
	slog.Debug("server stopped")
}
