package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const loggerContextKey contextKey = "logger"

var baseLogger = newLogger()

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(os.Getenv("LOG_LEVEL")),
	})
	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Init() *slog.Logger {
	slog.SetDefault(baseLogger)
	return baseLogger
}

func Default() *slog.Logger {
	return baseLogger
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Default()
	}
	if logger, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok && logger != nil {
		return logger
	}
	return Default()
}
