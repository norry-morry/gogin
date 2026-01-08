package log

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

var logKey ctxKey

// IntoLogger は context にロガーを埋め込みます。
func IntoLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

// FromLogger は context からロガーを取り出します（なければ no-op ロガー）。
func FromLogger(ctx context.Context) *slog.Logger {
	if v := ctx.Value(logKey); v != nil {
		if l, ok := v.(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return slog.Default()
}
