package logger

import (
	"log/slog"
	"os"
)

// Logger wraps slog.Logger
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance
func New(level string) *Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Fatal logs an error message and exits the program
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
