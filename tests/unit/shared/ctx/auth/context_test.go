// Package auth_test は context ベースの認証クレームヘルパ
// (internal/shared/ctx/auth/context.go) のユニットテストを提供します。
package auth_test

import (
	"context"
	"testing"

	authctx "resume/internal/shared/ctx/auth"
)

//
// ---------- With / From の基本動作 ----------
//

// TestWithAndFrom_Basic は With で設定した Claims を From で取得できることを確認します。
func TestWithAndFrom_Basic(t *testing.T) {
	t.Parallel()

	base := context.Background()
	claims := authctx.Claims{
		UID:   "uid-123",
		Email: "user@example.com",
	}

	ctx := authctx.With(base, claims)

	got, ok := authctx.From(ctx)
	if !ok {
		t.Fatalf("expected ok=true, got=false")
	}

	if got.UID != claims.UID || got.Email != claims.Email {
		t.Errorf("expected %+v, got %+v", claims, got)
	}
}

//
// ---------- From: 未設定 / 型不一致 ----------
//

// TestFrom_NotSet は Claims が設定されていない場合に ok=false になることを確認します。
func TestFrom_NotSet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	got, ok := authctx.From(ctx)
	if ok {
		t.Fatalf("expected ok=false when no claims set, got=true")
	}

	// got はゼロ値の Claims になっていること
	if got.UID != "" || got.Email != "" {
		t.Errorf("expected zero Claims, got %+v", got)
	}
}

// TestFrom_WrongType は 同じ Key ではない値が入っている場合に ok=false になることを確認します。
func TestFrom_WrongType(t *testing.T) {
	t.Parallel()

	type otherKey struct{}

	// 全く別のキーに値を入れておく
	ctx := context.WithValue(context.Background(), otherKey{}, authctx.Claims{
		UID:   "uid-xxx",
		Email: "wrong@example.com",
	})

	got, ok := authctx.From(ctx)
	if ok {
		t.Fatalf("expected ok=false for wrong key/type, got=true")
	}
	if got.UID != "" || got.Email != "" {
		t.Errorf("expected zero Claims on mismatch, got %+v", got)
	}
}
