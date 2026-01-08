// Package middleware presenter/http/camel_snake_codec.go
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// ContextKeySkipCodec コンテキストに true をセットすると、このミドルウェアは処理をスキップします。
//
//	c.Set(ContextKeySkipCodec, true)
const ContextKeySkipCodec = "skipCamelSnakeCodec"

// CamelSnakeOptions は camelCase/snake_case の相互変換ミドルウェアの動作を制御します。
type CamelSnakeOptions struct {
	EnableRequestDecode  bool // Request: camelCase -> snake_case
	EnableResponseEncode bool // Response: snake_case -> camelCase
}

func defaultOptions() CamelSnakeOptions {
	return CamelSnakeOptions{
		EnableRequestDecode:  true,
		EnableResponseEncode: true,
	}
}

// CamelSnakeCodec は JSON のキーをリクエストで snake_case に、レスポンスで camelCase に変換する Gin ミドルウェアを返します。
// オプションで各方向の変換を有効/無効にできます。
func CamelSnakeCodec(opts ...CamelSnakeOptions) gin.HandlerFunc {
	opt := defaultOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(c *gin.Context) {
		if skip, ok := c.Get(ContextKeySkipCodec); ok {
			if b, ok := skip.(bool); ok && b {
				c.Next()
				return
			}
		}

		// ---- Query: camelCase -> snake_case ----
		q := c.Request.URL.Query()
		changed := false
		for key, vals := range q {
			snake := camelToSnake(key)
			if snake != key {
				if _, exists := q[snake]; !exists {
					q[snake] = vals
					changed = true
				}
			}
		}
		if changed {
			c.Request.URL.RawQuery = q.Encode()
		}

		// ---- Request: camelCase -> snake_case ----
		if opt.EnableRequestDecode && isJSONContentType(c.Request.Header.Get("Content-Type")) &&
			c.Request.Body != nil && c.Request.ContentLength != 0 && c.Request.Method != http.MethodGet {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bytes.TrimSpace(bodyBytes)) > 0 {
				var any interface{}
				if json.Unmarshal(bodyBytes, &any) == nil {
					any = convertKeysCamelToSnake(any)
					if patched, err := json.Marshal(any); err == nil {
						c.Request.Body = io.NopCloser(bytes.NewReader(patched))
						c.Request.ContentLength = int64(len(patched))
						c.Request.Header.Del("Transfer-Encoding")
					} else {
						// 失敗時は元ボディを戻す
						c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
					}
				} else {
					// JSONでなければ元ボディを戻す
					c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
			}
		}

		// ---- Response: snake_case -> camelCase ----
		bw := &bufferedWriter{ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		status := bw.Status()
		if status == http.StatusNoContent || status == http.StatusNotModified || c.Request.Method == http.MethodHead {
			return
		}

		ct := bw.Header().Get("Content-Type")
		ce := bw.Header().Get("Content-Encoding")
		if bw.buf.Len() == 0 || !isJSONContentType(ct) || ce != "" || !opt.EnableResponseEncode {
			if bw.buf.Len() > 0 {
				writeRaw(bw, bw.buf.Bytes())
			}
			return
		}

		original := bw.buf.Bytes()
		var any interface{}
		if err := json.Unmarshal(original, &any); err != nil {
			// JSONじゃなければそのまま
			writeRaw(bw, original)
			return
		}
		any = convertKeysSnakeToCamel(any)
		modified, err := json.Marshal(any)
		if err != nil {
			// 変換失敗時はそのまま
			writeRaw(bw, original)
			return
		}

		bw.Header().Set("Content-Length", strconv.Itoa(len(modified)))
		bw.Header().Del("Transfer-Encoding")
		if _, err := bw.ResponseWriter.Write(modified); err != nil {
			// ここで失敗してもやれることは少ないのでログのみ
			log.Printf("camel_snake_codec: write modified response failed: %v", err)
		}
	}
}

// ---- Response Writer ラッパ ----

type bufferedWriter struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bufferedWriter) Write(b []byte) (int, error) {
	// ここではバッファに溜めるだけ。実書き込みはミドルウェアの後段で一度だけ行う。
	return w.buf.Write(b)
}

// 元の ResponseWriter に raw で書く（バッファは使わない）
func writeRaw(bw *bufferedWriter, b []byte) {
	if len(b) == 0 {
		return
	}
	bw.Header().Set("Content-Length", strconv.Itoa(len(b)))
	bw.Header().Del("Transfer-Encoding")
	if _, err := bw.ResponseWriter.Write(b); err != nil {
		log.Printf("camel_snake_codec: write raw response failed: %v", err)
	}
}

// ---- ユーティリティ ----

func isJSONContentType(ct string) bool {
	if ct == "" {
		return false
	}
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(ct)), "application/json")
}

func convertKeysSnakeToCamel(v interface{}) interface{} {
	switch vv := v.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{}, len(vv))
		for k, val := range vv {
			res[snakeToCamel(k)] = convertKeysSnakeToCamel(val)
		}
		return res
	case []interface{}:
		for i := range vv {
			vv[i] = convertKeysSnakeToCamel(vv[i])
		}
		return vv
	default:
		return v
	}
}

func convertKeysCamelToSnake(v interface{}) interface{} {
	switch vv := v.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{}, len(vv))
		for k, val := range vv {
			res[camelToSnake(k)] = convertKeysCamelToSnake(val)
		}
		return res
	case []interface{}:
		for i := range vv {
			vv[i] = convertKeysCamelToSnake(vv[i])
		}
		return vv
	default:
		return v
	}
}

// lowerCamel に揃える。Builder の戻り値を必ずチェックする。
func snakeToCamel(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	upperNext := false
	for i, r := range s {
		if r == '_' || r == '-' {
			upperNext = true
			continue
		}
		if i == 0 {
			if ok := safeWriteRune(&b, unicode.ToLower(r)); !ok {
				return s
			}
			continue
		}
		if upperNext {
			if ok := safeWriteRune(&b, unicode.ToUpper(r)); !ok {
				return s
			}
			upperNext = false
		} else {
			if ok := safeWriteRune(&b, r); !ok {
				return s
			}
		}
	}
	return b.String()
}

func camelToSnake(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				if ok := safeWriteByte(&b, '_'); !ok {
					return s
				}
			}
			if ok := safeWriteRune(&b, unicode.ToLower(r)); !ok {
				return s
			}
		} else if r == '-' {
			if ok := safeWriteByte(&b, '_'); !ok {
				return s
			}
		} else {
			if ok := safeWriteRune(&b, r); !ok {
				return s
			}
		}
	}
	return b.String()
}

// ---- Builder 安全書き込みヘルパ ----

func safeWriteRune(b *strings.Builder, r rune) bool {
	_, err := b.WriteRune(r)
	return err == nil
}

func safeWriteByte(b *strings.Builder, by byte) bool {
	if err := b.WriteByte(by); err != nil {
		return false
	}
	return true
}
