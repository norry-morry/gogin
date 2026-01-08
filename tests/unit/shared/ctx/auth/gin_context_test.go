// Package auth_test は Gin の Context 上で認証情報を扱うヘルパ
// (internal/shared/ctx/auth/gin_context.go) のユニットテストを提供します。
package auth_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	authctx "resume/internal/shared/ctx/auth"
)

//
// ---------- ヘルパ ----------
//

// newTestContext はテスト用の *gin.Context を生成します。
func newTestContext(t *testing.T) *gin.Context {
	t.Helper()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}

//
// ---------- UserID / SetUserID ----------
//

// TestUserID_SetAndGet は SetUserID で設定した値を UserID で取得できることを確認します。
func TestUserID_SetAndGet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	authctx.SetUserID(c, 42)

	got, ok := authctx.UserID(c)
	if !ok {
		t.Fatalf("expected ok=true, got=false")
	}
	if got != 42 {
		t.Errorf("expected 42, got=%d", got)
	}
}

// TestUserID_NotSet は 何も設定していない場合に (0, false) が返ることを確認します。
func TestUserID_NotSet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	got, ok := authctx.UserID(c)
	if ok {
		t.Fatalf("expected ok=false, got=true")
	}
	if got != 0 {
		t.Errorf("expected id=0, got=%d", got)
	}
}

// TestUserID_WrongType は 型が異なる値が入っている場合に ok=false になることを確認します。
func TestUserID_WrongType(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	// 誤って string などが入っているケース
	c.Set("userID", "not-uint64")

	got, ok := authctx.UserID(c)
	if ok {
		t.Fatalf("expected ok=false for wrong type, got=true")
	}
	if got != 0 {
		t.Errorf("expected id=0 on type mismatch, got=%d", got)
	}
}

//
// ---------- FirebaseUID / UID / SetFirebaseUID ----------
//

// TestFirebaseUID_SetAndGet は SetFirebaseUID で設定した値を FirebaseUID で取得できることを確認します。
func TestFirebaseUID_SetAndGet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	authctx.SetFirebaseUID(c, "firebase-uid-123")

	got, ok := authctx.FirebaseUID(c)
	if !ok {
		t.Fatalf("expected ok=true, got=false")
	}
	if got != "firebase-uid-123" {
		t.Errorf("expected firebase-uid-123, got=%s", got)
	}
}

// TestUID_CompatKey は SetFirebaseUID が UID 互換キーにも設定していることを確認します。
func TestUID_CompatKey(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	authctx.SetFirebaseUID(c, "uid-compat-456")

	got, ok := authctx.UID(c)
	if !ok {
		t.Fatalf("expected ok=true, got=false")
	}
	if got != "uid-compat-456" {
		t.Errorf("expected uid-compat-456, got=%s", got)
	}
}

// TestFirebaseUID_NotSet は FirebaseUID が未設定の場合に ("", false) が返ることを確認します。
func TestFirebaseUID_NotSet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	got, ok := authctx.FirebaseUID(c)
	if ok {
		t.Fatalf("expected ok=false, got=true")
	}
	if got != "" {
		t.Errorf("expected empty string, got=%s", got)
	}
}

// TestUID_NotSet は UID が未設定の場合に ("", false) が返ることを確認します。
func TestUID_NotSet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	got, ok := authctx.UID(c)
	if ok {
		t.Fatalf("expected ok=false, got=true")
	}
	if got != "" {
		t.Errorf("expected empty string, got=%s", got)
	}
}

//
// ---------- Email / SetEmail ----------
//

// TestEmail_SetAndGet は SetEmail で設定したメールアドレスを Email で取得できることを確認します。
func TestEmail_SetAndGet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	addr := "user@example.com"
	authctx.SetEmail(c, &addr)

	got, ok := authctx.Email(c)
	if !ok {
		t.Fatalf("expected ok=true, got=false")
	}
	if got == nil {
		t.Fatalf("expected non-nil email pointer")
	}
	if *got != addr {
		t.Errorf("expected %s, got=%s", addr, *got)
	}
}

// TestEmail_SetNil は SetEmail に nil を渡した場合に (nil, true) が返ることを確認します。
func TestEmail_SetNil(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	authctx.SetEmail(c, nil)

	got, ok := authctx.Email(c)
	if !ok {
		t.Fatalf("expected ok=true when key is set to nil, got=false")
	}
	if got != nil {
		t.Errorf("expected nil email pointer, got=%v", got)
	}
}

// TestEmail_NotSet は Email が未設定の場合に (nil, false) が返ることを確認します。
func TestEmail_NotSet(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	got, ok := authctx.Email(c)
	if ok {
		t.Fatalf("expected ok=false when email is not set, got=true")
	}
	if got != nil {
		t.Errorf("expected nil email pointer when not set, got=%v", got)
	}
}

// TestEmail_WrongType は 異なる型が入っている場合に ok=false になることを確認します。
func TestEmail_WrongType(t *testing.T) {
	t.Parallel()

	c := newTestContext(t)

	// 誤って string 型そのものが入っているケース
	c.Set("email", "not-pointer")

	got, ok := authctx.Email(c)
	if ok {
		t.Fatalf("expected ok=false for wrong type, got=true")
	}
	if got != nil {
		t.Errorf("expected nil email pointer on type mismatch, got=%v", got)
	}
}
