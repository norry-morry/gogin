// Package requestid_test は requestid ユーティリティのユニットテスト
// (internal/shared/requestid/requestid.go) を提供します。
package requestid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"resume/internal/shared/requestid"
)

//
// ---------- Ensure（既存値 → ヘッダ → 生成） ----------
//

// TestEnsure_ContextHasValue は Context に既存の request_id があればそれを返すことを確認する。
func TestEnsure_ContextHasValue(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("request_id", "ctx_123")

	got := requestid.Ensure(c)

	if got != "ctx_123" {
		t.Errorf("expected ctx_123, got=%s", got)
	}

	// レスポンスヘッダにもセットされる
	if h := w.Header().Get(requestid.HeaderCanonical); h != "ctx_123" {
		t.Errorf("expected response header=%s, got=%s", "ctx_123", h)
	}
}

// TestEnsure_HeaderProvided は X-Request-ID があればそれを採用することを確認する。
func TestEnsure_HeaderProvided(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "hdr_456")
	c.Request = req

	got := requestid.Ensure(c)

	if got != "hdr_456" {
		t.Errorf("expected hdr_456, got=%s", got)
	}

	if h := w.Header().Get(requestid.HeaderCanonical); h != "hdr_456" {
		t.Errorf("expected header=%s, got=%s", "hdr_456", h)
	}
}

// TestEnsure_HeaderLowercaseVariant は X-Request-Id（小文字違い）でも拾えることを確認する。
func TestEnsure_HeaderLowercaseVariant(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "mixcase_789")
	c.Request = req

	got := requestid.Ensure(c)

	if got != "mixcase_789" {
		t.Errorf("expected mixcase_789, got=%s", got)
	}

	if h := w.Header().Get(requestid.HeaderCanonical); h != "mixcase_789" {
		t.Errorf("expected header=%s, got=%s", "mixcase_789", h)
	}
}

// TestEnsure_GenerateNew は ContextにもHeaderにも無い場合、新規生成されることを確認する。
func TestEnsure_GenerateNew(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// ★ これを追加
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	got := requestid.Ensure(c)

	if !strings.HasPrefix(got, "backend_") {
		t.Errorf("expected prefix backend_, got=%s", got)
	}

	if h := w.Header().Get(requestid.HeaderCanonical); h != got {
		t.Errorf("header mismatch: expected %s, got %s", got, h)
	}
}

//
// ---------- Get（Ensure 後に取得できる） ----------
//

// TestGet_ReturnsValue は Ensure が設定した値を取得できることを確認。
func TestGet_ReturnsValue(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// ★ これを追加
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	requestid.Ensure(c)
	got := requestid.Get(c)

	if got == "" {
		t.Errorf("expected non-empty request id")
	}
}

// TestGet_NoValue は 未設定なら空文字を返すことを確認。
func TestGet_NoValue(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if got := requestid.Get(c); got != "" {
		t.Errorf("expected empty, got=%s", got)
	}
}
