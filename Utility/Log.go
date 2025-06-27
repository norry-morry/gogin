package Utility

import (
	"context"
	"log/slog"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if e := h.Handle(ctx, r); e != nil && err == nil {
				err = e
			}
		}
	}
	return err
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var wrapped []slog.Handler
	for _, h := range m.handlers {
		wrapped = append(wrapped, h.WithAttrs(attrs))
	}
	return &MultiHandler{handlers: wrapped}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	var wrapped []slog.Handler
	for _, h := range m.handlers {
		wrapped = append(wrapped, h.WithGroup(name))
	}
	return &MultiHandler{handlers: wrapped}
}
