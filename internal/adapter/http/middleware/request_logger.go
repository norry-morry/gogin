// Package middleware はリクエストのロギングなどのミドルウェアを提供します。
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"resume/internal/shared/requestid"
	"resume/internal/shared/util"
)

// RequestLoggerOptions はリクエストロガーミドルウェアの動作を制御します。
type RequestLoggerOptions struct {
	// true を返したリクエストはロギングしない（/health など）
	Skipper func(*gin.Context) bool
	// 追加で残したいヘッダ（値は req.headers にマップ化）
	IncludeHeaders []string
	// 追加：任意の Masker。nil の場合は util.GlobalMasker() を使う。
	Masker util.Masker
}

//const ctxKeyRawBody = "request.rawBody" // 後段で見たいとき用に保存

// RequestLogger は Gin のリクエスト/レスポンス情報を構造化ログで出力するミドルウェアを返します。
func RequestLogger(l *slog.Logger, opt *RequestLoggerOptions) gin.HandlerFunc {
	// nil セーフな初期化
	var (
		skipper       func(*gin.Context) bool
		includeHeader = map[string]struct{}{}
		masker        util.Masker
	)
	if opt != nil && opt.Skipper != nil {
		skipper = opt.Skipper
	}
	if opt != nil && len(opt.IncludeHeaders) > 0 {
		for _, h := range opt.IncludeHeaders {
			includeHeader[strings.ToLower(h)] = struct{}{}
		}
	}
	// マスカーは指定がなければグローバルを使用
	if opt != nil && opt.Masker != nil {
		masker = opt.Masker
	} else {
		masker = util.GlobalMasker()
	}

	normalize := func(h string) string { return strings.ToLower(h) }

	return func(c *gin.Context) {
		// スキップ条件
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}

		start := time.Now()

		// Request 基本情報
		reqID := requestid.Ensure(c)
		method := c.Request.Method
		path := c.FullPath()
		if path == "" { // ルート未マッチ時は生パス
			path = c.Request.URL.Path
		}
		host := c.Request.Host
		proto := c.Request.Proto
		ip := c.ClientIP()
		ua := c.Request.UserAgent()

		// ---- クエリ（マスク対象）----
		qmap := make(map[string]any, len(c.Request.URL.Query()))
		for k, vs := range c.Request.URL.Query() {
			if len(vs) > 0 {
				qmap[k] = vs[0]
			}
		}
		maskedQuery := masker.MaskMapShallow(qmap)

		// ---- ヘッダ（IncludeHeaders のみ・マスク対象）----
		hmap := map[string]any{}
		for k, v := range c.Request.Header {
			if _, ok := includeHeader[normalize(k)]; ok && len(v) > 0 {
				hmap[k] = masker.MaskByKey(k, v[0])
			}
		}

		// ---- ボディ（JSON優先で判定→再帰マスク。非JSONはそのまま）----
		var bodyLogged any
		if c.Request.Body != nil && c.Request.ContentLength != 0 {
			if b, err := io.ReadAll(c.Request.Body); err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(b)) // 差し戻し

				var anyJSON any
				if len(b) > 0 && json.Unmarshal(b, &anyJSON) == nil {
					// ✅ JSON として読めたら再帰マスク
					masked := masker.MaskAnyRecursive("", anyJSON)
					if enc, err := json.Marshal(masked); err == nil {
						bodyLogged = string(enc)
					} else {
						bodyLogged = string(b) // フォールバック
					}
				} else {
					// 非JSONはそのまま（必要ならここで独自マスク）
					bodyLogged = string(b)
				}
			}
		}

		// 実処理
		c.Next()

		// Response 側
		status := c.Writer.Status()
		size := c.Writer.Size() // -1 の場合あり
		latMs := time.Since(start).Milliseconds()

		// Gin のエラー（c.Error されたもの）
		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		attrs := []any{
			"request_id", reqID,
			slog.Group("req",
				"method", method,
				"path", path,
				"host", host,
				"proto", proto,
				"ip", ip,
				"user_agent", ua,
				"query", maskedQuery, // ← マスク済み
				"headers", hmap, // ← マスク済み
				"content_length", c.Request.ContentLength,
				"body", bodyLogged, // JSON はマスク済み文字列
			),
			slog.Group("res",
				"status", status,
				"size", size,
			),
			"duration_ms", latMs,
		}

		if errMsg != "" {
			attrs = append(attrs, "error", errMsg)
			l.Error("http.request", attrs...)
			return
		}
		l.Info("http.request", attrs...)
	}
}
