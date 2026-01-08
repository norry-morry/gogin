// Package firebase provides factories and interfaces for Firebase Admin Auth.
package firebase

import (
	"context"

	fb "firebase.google.com/go/v4/auth"
)

// Auth は Firebase Authentication のラッパーインターフェースです。
// アクセストークンの検証や、Firebaseユーザー情報の取得など、
// 認証関連の操作を抽象化して提供します。
type Auth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*fb.Token, error)
	GetUser(ctx context.Context, uid string) (*fb.UserRecord, error)
}

type authWrapper struct{ c *fb.Client }

// NewAuth は *auth.Client をアプリ用の Auth インターフェースでラップする。
func NewAuth(c *fb.Client) Auth { return &authWrapper{c: c} }

// VerifyIDToken は、Firebase が発行した ID トークンの署名と有効期限を検証し、トークン内のクレーム情報を返します
func (w *authWrapper) VerifyIDToken(ctx context.Context, t string) (*fb.Token, error) {
	return w.c.VerifyIDToken(ctx, t)
}

// GetUser は、VerifyTokenで得られた UID を元にFirebaseのユーザーを取得します
func (w *authWrapper) GetUser(ctx context.Context, uid string) (*fb.UserRecord, error) {
	return w.c.GetUser(ctx, uid)
}
