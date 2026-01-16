package logging

import (
	"log/slog"
	"os"
	"strings"
)

var defaultLogger *slog.Logger

func init() {
	defaultLogger = New()
}

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	format := strings.ToLower(os.Getenv("LOG_FORMAT"))

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

func L() *slog.Logger {
	return defaultLogger
}
