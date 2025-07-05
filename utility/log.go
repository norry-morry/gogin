// Package utility は共通のユーティリティ関数やロギング処理を提供します。
package utility

import (
	"context"
	"log/slog"
)

type multiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler は複数の slog.Handler をまとめて同時に処理する Handler を返します。
func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
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

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var wrapped []slog.Handler
	for _, h := range m.handlers {
		wrapped = append(wrapped, h.WithAttrs(attrs))
	}
	return &multiHandler{handlers: wrapped}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	var wrapped []slog.Handler
	for _, h := range m.handlers {
		wrapped = append(wrapped, h.WithGroup(name))
	}
	return &multiHandler{handlers: wrapped}
}
