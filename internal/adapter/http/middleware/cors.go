package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig は CORS ミドルウェアの動作を制御します
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig は CORS レスポンスヘッダの設定を保持する構造体です。
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     envCSV("CORS_ALLOW_ORIGINS", []string{"http://resume.local", "http://www.resume.local", "http://localhost:3000", "http://localhost:5173"}),
		AllowMethods:     envCSV("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowHeaders:     envCSV("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Request-ID"}),
		ExposeHeaders:    envCSV("CORS_EXPOSE_HEADERS", []string{"X-Request-ID"}),
		AllowCredentials: envBool("CORS_ALLOW_CREDENTIALS", false),
		MaxAge:           envDuration("CORS_MAX_AGE", 12*time.Hour),
	}
}

// Cors NewCORS returns a Gin middleware that sets Cors headers.
func Cors(cfg CORSConfig) gin.HandlerFunc {
	allowAll := contains(cfg.AllowOrigins, "*")
	allowMethods := strings.Join(cfg.AllowMethods, ",")
	fallbackAllowHeaders := strings.Join(cfg.AllowHeaders, ",")
	expose := strings.Join(cfg.ExposeHeaders, ",")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Next()
			return
		}
		// 許可判定（完全一致 or "*"）
		originAllowed := allowAll || containsFold(cfg.AllowOrigins, origin)
		if !originAllowed {
			c.Next()
			return
		}

		// ここで“常時”CORSをセット（後で上書きされても WriteHeader/Write でもう一度付ける）
		setBase := func() {
			h := c.Writer.Header()
			// withCredentials=true の場合は "*" を使えない
			if cfg.AllowCredentials || !allowAll {
				h.Set("Access-Control-Allow-Origin", origin) // 単一値
			} else {
				h.Set("Access-Control-Allow-Origin", "*")
			}
			if cfg.AllowCredentials {
				h.Set("Access-Control-Allow-Credentials", "true")
			}
			if expose != "" {
				h.Set("Access-Control-Expose-Headers", expose)
			}
			addVary(c, "Origin")
		}
		setBase()

		// プリフライトはここで完結
		if c.Request.Method == http.MethodOptions &&
			c.GetHeader("Access-Control-Request-Method") != "" {

			reqHdr := c.GetHeader("Access-Control-Request-Headers")
			if reqHdr == "" {
				reqHdr = fallbackAllowHeaders
			}
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Methods", allowMethods)
			h.Set("Access-Control-Allow-Headers", reqHdr)
			if cfg.MaxAge > 0 {
				h.Set("Access-Control-Max-Age", strconv.Itoa(int(cfg.MaxAge/time.Second)))
			}
			addVary(c, "Access-Control-Request-Method")
			addVary(c, "Access-Control-Request-Headers")

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// --- ここがキモ：最終書き込み直前にもう一度 CORS を強制注入 ---
		origWriter := c.Writer
		c.Writer = &corsResponseWriter{
			ResponseWriter: origWriter,
			ensure:         setBase,
		}

		c.Next()
	}
}

// gin.ResponseWriter をラップし、WriteHeader/Write の直前で ensure() を呼ぶ
type corsResponseWriter struct {
	gin.ResponseWriter
	ensure func()
}

func (w *corsResponseWriter) WriteHeader(code int) {
	w.ensure()
	w.ResponseWriter.WriteHeader(code)
}
func (w *corsResponseWriter) Write(b []byte) (int, error) {
	// header 未送出のまま Write が先に呼ばれる場合に備えて
	w.ensure()
	return w.ResponseWriter.Write(b)
}

func addVary(c *gin.Context, v string) {
	const key = "Vary"
	cur := c.Writer.Header().Get(key)
	if cur == "" {
		c.Header(key, v)
		return
	}
	// 既に含まれていれば重複させない
	for _, part := range strings.Split(cur, ",") {
		if strings.EqualFold(strings.TrimSpace(part), v) {
			return
		}
	}
	c.Header(key, cur+", "+v)
}

func envCSV(key string, def []string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return def
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func envBool(key string, def bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return def
	}
	switch raw {
	case "1", "true", "t", "yes", "y":
		return true
	case "0", "false", "f", "no", "n":
		return false
	default:
		return def
	}
}

func envDuration(key string, def time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return def
	}
	// time.ParseDuration 形式（例: "2h", "30m"）
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	return def
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func containsFold(list []string, v string) bool {
	for _, s := range list {
		if strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(v)) {
			return true
		}
	}
	return false
}

//func formatSeconds(d time.Duration) string {
//	sec := int(d / time.Second)
//	return strconvItoa(sec)
//}

// ループ回避のため最小限
//func strconvItoa(i int) string {
//	// 依存を減らすために簡易実装でもOKだが、通常は strconv.Itoa を使う
//	// ここでは標準を使います
//	return strconv.Itoa(i)
//}
