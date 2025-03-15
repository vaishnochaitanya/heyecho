package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// SimpleHandler is a minimal handler that formats logs without prefixes
type SimpleHandler struct {
	w     io.Writer
	level slog.Level
}

// NewSimpleHandler creates a new simple log handler
func NewSimpleHandler(w io.Writer, level slog.Level) *SimpleHandler {
	return &SimpleHandler{
		w:     w,
		level: level,
	}
}

func (h *SimpleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *SimpleHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006/01/02 15:04:05")

	fmt.Fprintf(h.w, "[%s] %s: %s", timeStr, r.Level, r.Message)

	// Add any additional attributes
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value.Any())
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we ignore attrs in this implementation
	return h
}

func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we ignore groups in this implementation
	return h
}
