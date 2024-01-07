// Package log provides a small abstraction around slog to make it easier to use.
package log

import (
	"context"
	"log/slog"
)

var (
	Int    = slog.Int
	String = slog.String
	Bool   = slog.Bool
)

// Info logs a message with info level
func Info(ctx context.Context, message string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelInfo, message, attrs...)
}

// Error logs a message with error level
func Error(ctx context.Context, message string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelError, message, attrs...)
}

// Debug logs a message with debug level
func Debug(ctx context.Context, message string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelDebug, message, attrs...)
}
