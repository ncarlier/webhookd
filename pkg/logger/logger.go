package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	// HookOutputEnabled writes hook output into logs if true
	HookOutputEnabled = false
	// RequestOutputEnabled writes HTTP request into logs if true
	RequestOutputEnabled = false
)

// Configure logger
func Configure(format, level string) {
	logLevel := slog.LevelDebug
	switch level {
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	var logger *slog.Logger
	if format == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &opts))
	}

	slog.SetDefault(logger)
}

// LogIf writ log on condition
func LogIf(condition bool, level slog.Level, msg string, args ...any) {
	if condition {
		slog.Log(context.Background(), level, msg, args...)
	}
}
