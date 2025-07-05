// Package middleware はリクエストのロギングなどのミドルウェアを提供します。
package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type ctxKeyLogger struct{}

var loggerKey = ctxKeyLogger{}

// SetLoggerContext は ctx に構造化 logger を埋め込む
func SetLoggerContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// GetLogger は ctx から logger を取得。なければ slog.Default()
func GetLogger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok || l == nil {
		return slog.Default()
	}
	return l
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Request-ID を取得（なければタイムスタンプ）
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = time.Now().Format("20060102150405.000")
		}

		// 構造化ログ作成
		reqLogger := slog.Default().With(
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		// context に logger を埋め込む
		// context.WithValue を使う（Go 1.21 用）
		//ctx := context.WithValue(c.Request.Context(), loggerKey, reqLogger)
		//c.Request = c.Request.WithContext(ctx)
		ctx := SetLoggerContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logFunc := reqLogger.Info
		if status >= 500 {
			logFunc = reqLogger.Error
		}

		logFunc("completed request",
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)
	}
}
