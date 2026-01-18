package logging

import (
	"log/slog"
	"os"
	"strings"
)

var defaultLogger *slog.Logger
var verbose bool

func init() {
	defaultLogger = New()
}

func New() *slog.Logger {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Timestamps and source
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

func SetVerbose(v bool) {
	verbose = v
	defaultLogger = New()
}

func L() *slog.Logger {
	return defaultLogger
}
