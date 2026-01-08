package auth

import "context"

// Key は認証クレームを context に入れる際のキーです。
type Key struct{}

// Claims は認証済みユーザーのクレーム情報です。
type Claims struct {
	UID   string
	Email string
}

// With はクレームを context に詰めて返します。
func With(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, Key{}, c)
}

// From は context からクレームを取り出します。
func From(ctx context.Context) (Claims, bool) {
	v, ok := ctx.Value(Key{}).(Claims)
	return v, ok
}
