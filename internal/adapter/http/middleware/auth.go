package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"resume/internal/adapter/gateway/firebase"
	"resume/internal/adapter/http/handlerutil"
	"resume/internal/shared/apperr"
	"resume/internal/shared/ctx/auth"
)

// Bearer extracts the token from an Authorization header in "Bearer <token>" format.
// It returns the token and true if successful, otherwise an empty string and false.
func Bearer(h string) (string, bool) {
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// InjectAuth は Authorization ヘッダーの Bearer トークンを検証し、
// 検証に成功した場合は認証クレーム（UID と email）をリクエストコンテキストへ注入します。
func InjectAuth(fb firebase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, ok := Bearer(c.GetHeader("Authorization"))
		if !ok {
			c.Next()
			return
		}
		t, err := fb.VerifyIDToken(c.Request.Context(), tok)
		if err != nil {
			// 無効なら素通り or 401 を返すかは方針次第。ここでは素通りに。
			c.Next()
			return
		}
		email := ""
		if v, ok := t.Claims["email"].(string); ok {
			email = v
		}
		ctx := auth.With(c.Request.Context(), auth.Claims{
			UID: t.UID, Email: email,
		})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireAuth は Bearer トークンの存在と有効性を強制し、
// 失敗した場合は 401 を返して処理を中断します。成功時は認証クレームを注入します。
func RequireAuth(fb firebase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ★ プリフライトは通す
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		tok, ok := Bearer(c.GetHeader("Authorization"))
		if !ok {
			handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
			return
		}

		t, err := fb.VerifyIDToken(c.Request.Context(), tok)
		if err != nil {
			handlerutil.WriteError(
				c,
				apperr.Wrap(apperr.CodeUnauthorized, "invalid token", err, nil),
			)
			return
		}

		email := ""
		if v, ok := t.Claims["email"].(string); ok {
			email = v
		}
		ctx := auth.With(c.Request.Context(), auth.Claims{UID: t.UID, Email: email})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
