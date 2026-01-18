// Package requestid は、Gin のコンテキストと HTTP ヘッダ間で
// リクエストIDを取得・生成・伝播するユーティリティを提供します。
package requestid

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderCanonical は、レスポンスに設定するリクエストIDヘッダ名の正規表記です。
	// 入力側では "X-Request-ID" / "X-Request-Id" の双方を受け入れます。
	HeaderCanonical = "X-Request-ID" // 統一表記（出力はこれ）
	ctxKey          = "request_id"
)

var acceptedHeaders = []string{
	"X-Request-ID",
	"X-Request-Id",
	//"X-Trace-Id",
	//"X-Correlation-Id",
}

// Ensure は、既存のリクエストIDをヘッダまたはコンテキストから取得し、
// なければ新規生成してコンテキストとレスポンスヘッダへ設定して返します。
func Ensure(c *gin.Context) string {
	// 1) Contextに既にあればそれを使う（※ヘッダにも揃える）
	if rid, ok := c.Get(ctxKey); ok {
		if s, ok := rid.(string); ok && s != "" {
			return setAll(c, s)
		}
	}

	// 2) 受け入れヘッダから探す（大小無視）
	for _, h := range acceptedHeaders {
		if v := c.GetHeader(h); v != "" {
			return setAll(c, v)
		}
	}

	// 3) なければ生成
	rid := uuid.NewString()
	return setAll(c, "backend_"+rid)
}

func setAll(c *gin.Context, id string) string {
	id = strings.TrimSpace(id)
	c.Set(ctxKey, id)                          // ハンドラ用
	c.Writer.Header().Set(HeaderCanonical, id) // レスポンス用
	return id
}

// Get は Ensure が実行済みの場合にリクエストIDを返します。未設定なら空文字を返します。
func Get(c *gin.Context) string {
	if rid, ok := c.Get(ctxKey); ok {
		if s, ok := rid.(string); ok && s != "" {
			return s
		}
	}
	return ""
}
