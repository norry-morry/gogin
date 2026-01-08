// Package controller はエンドポイント処理を担当します。
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	fba "resume/internal/adapter/gateway/firebase"
	"resume/internal/adapter/http/handlerutil"
	"resume/internal/adapter/http/middleware"
	"resume/internal/shared/apperr"
	"resume/internal/shared/util"
	ucauth "resume/internal/usecase/auth"
)

// AuthHandler は認証用のハンドラです。
type AuthHandler struct {
	fb fba.Auth
	uc ucauth.Usecase
}

// NewAuthHandler は AuthHandler の新しいインスタンスを生成して返します。
// 認証関連の HTTP リクエストを処理するためのハンドラを初期化します。
func NewAuthHandler(fb fba.Auth, uc ucauth.Usecase) *AuthHandler {
	return &AuthHandler{fb: fb, uc: uc}
}

// UserLogin はユーザーのログイン情報を返します
// コンテキストから、idTokenを取得しユーザー情報の
func (h *AuthHandler) UserLogin(c *gin.Context) {
	idToken, ok := middleware.Bearer(c.GetHeader("Authorization"))
	if !ok || idToken == "" {
		handlerutil.WriteError(c,
			apperr.New(apperr.CodeBadRequest, "missing or invalid Authorization header", nil))
		return
	}

	// Verify token
	tok, err := h.fb.VerifyIDToken(c.Request.Context(), idToken)
	if err != nil {
		handlerutil.WriteError(c,
			apperr.Wrap(apperr.CodeUnauthorized, "invalid token", err, nil))
		return
	}

	// Fetch Firebase User
	ur, err := h.fb.GetUser(c.Request.Context(), tok.UID)
	if err != nil {
		handlerutil.WriteError(c,
			apperr.Wrap(apperr.CodeUnauthorized, "failed to fetch firebase user", err, nil))
		return
	}

	// Build Providers From UserRecord
	providers := make([]ucauth.ProviderInfo, 0, len(ur.ProviderUserInfo))
	for _, p := range ur.ProviderUserInfo {
		// p.ProviderID: "google.com" | "github.com" | "password" | ...
		providers = append(providers, ucauth.ProviderInfo{
			Provider:            util.SafeStr(p.ProviderID),
			ProviderUserID:      util.SafeStr(p.UID),
			ProviderDisplayName: util.FallbackStr(p.DisplayName, ur.DisplayName, "unknown"),
			EmailAtSignup:       util.OptPtr(p.Email),
		})
	}

	// Build input DTO (users/auth_identitiesに寄せた最小構成)
	in := ucauth.UpsertUserFromFirebaseInput{
		IDToken:       idToken,
		UID:           ur.UID,
		Email:         util.OptPtr(ur.Email),
		EmailVerified: ur.EmailVerified,
		DisplayName:   util.OptPtr(ur.DisplayName),
		PhotoURL:      util.OptPtr(ur.PhotoURL),
		Providers:     providers,
	}

	// Call UseCase (内部でlast_login_at を now に更新 & Upsert)
	out, err := h.uc.UpsertUserFromFirebase(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err) // usecase 側が apperr.Error を返すのでそのまま渡す
		return
	}

	c.JSON(http.StatusOK, out)
}
