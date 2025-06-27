package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type ctxKeyLogger struct{}

var LoggerKey = ctxKeyLogger{}

func RequestLoggerMiddleware(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = time.Now().Format("20060102150405.000")
		}

		// リクエストスコープの logger を作成
		reqLogger := baseLogger.With(
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
		)

		// context.WithValue を使う（Go 1.21 用）
		ctx := context.WithValue(c.Request.Context(), LoggerKey, reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		reqLogger.Info("completed request",
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)
	}
}
