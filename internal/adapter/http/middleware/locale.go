// Package middleware は、HTTPリクエストからロケールを交渉・取得するミドルウェアを提供します。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// LocaleSource は利用可能ロケールコードの列挙だけを要求する最小インターフェイスです。
// 旧実装の Store がこのメソッドを持っていれば、そのまま渡せます。
type LocaleSource interface {
	AvailableLocaleCodes() []string
}

// Opts はロケールネゴシエーションの設定です（旧Optsに準拠）。
type Opts struct {
	Default   string       // 既定言語（例: "ja"）
	Allow     []string     // 明示指定。空なら Store から自動検出
	Query     string       // 例: "lang"（空なら無効）
	Header    string       // 例: "Accept-Language"
	Cookie    string       // 例: "lang"（空なら無効）
	SetHeader bool         // Content-Language を付与するか（既定: true）
	Store     LocaleSource // 自動検出用（Allow が空のとき推奨）
}

// CtxLocaleKey は Gin Context に格納するキー名
const CtxLocaleKey = "lang"

// Middleware は、Header > Query > Cookie > Default の順にロケールを決定し、
// Context にセットします（必要なら Response Header へ Content-Language も付与）。
func Middleware(o Opts) gin.HandlerFunc {
	if o.Header == "" {
		o.Header = "Accept-Language"
	}
	// 旧実装同様、明示指定されなければ true
	if !o.SetHeader {
		o.SetHeader = true
	}
	return func(c *gin.Context) {
		loc := Negotiate(c.Request, o)
		if o.SetHeader {
			c.Writer.Header().Set("Content-Language", loc)
		}
		c.Set(CtxLocaleKey, loc)
		c.Next()
	}
}

// From は Context からロケールを取り出すヘルパーです。
func From(c *gin.Context) string {
	if v, ok := c.Get(CtxLocaleKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Negotiate は HTTP リクエストからロケールを決定して返します。
// 優先順位：Query > Header > Cookie > Default
func Negotiate(r *http.Request, o Opts) string {
	inAllow := func(code string) bool {
		code = strings.ToLower(strings.TrimSpace(code))
		if code == "" {
			return false
		}
		// Allow 未指定なら Store から検出
		allow := o.Allow
		if len(allow) == 0 && o.Store != nil {
			allow = o.Store.AvailableLocaleCodes()
		}
		if len(allow) == 0 {
			return true // 制限なし
		}
		for _, a := range allow {
			if strings.EqualFold(a, code) {
				return true
			}
		}
		return false
	}

	// 1) Query
	if qk := strings.TrimSpace(o.Query); qk != "" {
		if v := strings.TrimSpace(r.URL.Query().Get(qk)); v != "" && inAllow(v) {
			return strings.ToLower(v)
		}
	}

	// 2) Header（先頭、地域は落として base を試す）
	if o.Header != "" {
		al := r.Header.Get(o.Header)
		if al != "" {
			for _, part := range strings.Split(al, ",") {
				raw := strings.TrimSpace(strings.Split(part, ";")[0])
				if raw == "" {
					continue
				}
				raw = strings.ToLower(raw)
				// そのまま
				if inAllow(raw) {
					return raw
				}
				// サブタグ切り落とし（例: en-US -> en）
				if i := strings.IndexByte(raw, '-'); i > 0 {
					base := raw[:i]
					if inAllow(base) {
						return base
					}
				}
			}
		}
	}

	// 3) Cookie
	if ck := strings.TrimSpace(o.Cookie); ck != "" {
		if c, err := r.Cookie(ck); err == nil && c != nil {
			if v := strings.TrimSpace(c.Value); v != "" && inAllow(v) {
				return strings.ToLower(v)
			}
		}
	}

	// 4) Default（最終フォールバック）
	return strings.ToLower(strings.TrimSpace(o.Default))
}

// Setter は Cookie を使ってロケールを保存するためのハンドラです。
// 例: GET /lang?to=en
func Setter(o Opts) gin.HandlerFunc {
	return func(c *gin.Context) {
		to := strings.TrimSpace(c.Query("to"))
		if to == "" {
			to = o.Default
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     o.Cookie,
			Value:    to,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400 * 365,
		})
		c.JSON(http.StatusOK, gin.H{"ok": true, "locale": to})
	}
}
