package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"resume/internal/adapter/http/handlerutil"
	"resume/internal/domain/repository"
	"resume/internal/shared/apperr"
	"resume/internal/shared/ctx/auth"
)

// ResolveUser は、Firebase UID からユーザー情報を解決するミドルウェアです。
type ResolveUser struct {
	userRepo repository.UserRepository
}

// NewResolveUser は、ユーザー解決ミドルウェアを初期化して返します。
func NewResolveUser(userRepo repository.UserRepository) *ResolveUser {
	return &ResolveUser{userRepo: userRepo}
}

// InjectResolveUser --- 1) 認証“任意”版: 認証が入っていれば user_id を解決・セット。無ければ素通り。
func (m *ResolveUser) InjectResolveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// OPTIONS は素通り（CORSプリフライト対策。RequireAuth と合わせる）
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := auth.From(c.Request.Context())
		if !ok || claims.UID == "" {
			// 認証なし → そのまま素通り
			c.Next()
			return
		}

		u, err := m.userRepo.FindByUID(c.Request.Context(), claims.UID)
		if err != nil {
			handlerutil.WriteError(c, apperr.Wrap(apperr.CodeInternal, "failed to load user", err, nil))
			return
		}
		if u == nil || u.ID == 0 {
			// 「任意」版なので未登録でも素通り（必要ならここで消すこともできる）
			c.Next()
			return
		}

		auth.SetUserID(c, uint64(u.ID))
		c.Next()
	}
}

// RequireResolveUser --- 2) 認証“必須”版: ユーザー未登録や解決失敗ならエラーで打ち切り。
func (m *ResolveUser) RequireResolveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// OPTIONS は素通り
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := auth.From(c.Request.Context())
		if !ok || claims.UID == "" {
			handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
			return
		}

		u, err := m.userRepo.FindByUID(c.Request.Context(), claims.UID)
		if err != nil {
			handlerutil.WriteError(c, apperr.Wrap(apperr.CodeInternal, "failed to load user", err, nil))
			return
		}
		if u == nil || u.ID == 0 {
			// チーム方針で 401/404/412 のどれかに揃えてOK。ここは 401 に寄せています。
			handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "user not registered", nil))
			return
		}

		auth.SetUserID(c, uint64(u.ID))
		c.Next()
	}
}
