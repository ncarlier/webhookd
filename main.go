package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
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

const envPrefix = "WHD"

func main() {
	conf := &config.Config{}
	configflag.Bind(conf, envPrefix)

	flag.Parse()

	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	if conf.Hook.LogDir == "" {
		conf.Hook.LogDir = os.TempDir()
	}

	if err := conf.Validate(); err != nil {
		log.Fatal("invalid configuration:", err)
	}

	logger.Configure(conf.Log.Format, conf.Log.Level)
	logger.HookOutputEnabled = slices.Contains(conf.Log.Modules, "hook")
	logger.RequestOutputEnabled = slices.Contains(conf.Log.Modules, "http")

	conf.ManageDeprecatedFlags(envPrefix)

	slog.Debug("starting webhookd server...")

	srv := server.NewServer(conf)

	// Configure notification
	if err := notification.Init(conf.Notification.URI); err != nil {
		slog.Error("unable to create notification channel", "err", err)
	}

	// Start the dispatcher.
	slog.Debug("starting the dispatcher...", "workers", conf.Hook.Workers)
	worker.StartDispatcher(conf.Hook.Workers)

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
