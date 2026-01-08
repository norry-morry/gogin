// Package apperr_test は application error 型
// (internal/shared/apperr/error.go) のユニットテストを提供します。
package apperr_test

import (
	"errors"
	"testing"

	"resume/internal/shared/apperr"
)

//
// ---------- Error.Error ----------
//

// TestError_Error_WithMessage は Message がある場合に "CODE: msg" 形式になることを確認します。
func TestError_Error_WithMessage(t *testing.T) {
	t.Parallel()

	e := &apperr.Error{
		Code:    apperr.CodeBadRequest,
		Message: "invalid payload",
	}

	got := e.Error()
	want := "BAD_REQUEST: invalid payload"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// TestError_Error_WithoutMessage は Message が空の場合に Code のみ文字列化されることを確認します。
func TestError_Error_WithoutMessage(t *testing.T) {
	t.Parallel()

	e := &apperr.Error{
		Code: apperr.CodeInternal,
	}

	got := e.Error()
	want := "INTERNAL"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

//
// ---------- New / Wrap / Unwrap ----------
//

// TestNew_ConstructsError は New が Cause=nil の Error を返すことを確認します。
func TestNew_ConstructsError(t *testing.T) {
	t.Parallel()

	details := map[string]any{"field": "name"}
	e := apperr.New(apperr.CodeUnprocessable, "validation error", details)

	if e.Code != apperr.CodeUnprocessable {
		t.Errorf("expected CodeUnprocessable, got %s", e.Code)
	}
	if e.Message != "validation error" {
		t.Errorf("unexpected message: %s", e.Message)
	}
	if e.Details["field"] != "name" {
		t.Errorf("unexpected details: %#v", e.Details)
	}
	if e.Cause != nil {
		t.Errorf("expected Cause=nil, got=%v", e.Cause)
	}
}

// TestWrap_ConstructsErrorWithCause は Wrap が Cause 付き Error を返すことを確認します。
func TestWrap_ConstructsErrorWithCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("db error")
	details := map[string]any{"retryable": true}

	e := apperr.Wrap(apperr.CodeInternal, "failed to save", cause, details)

	if e.Code != apperr.CodeInternal {
		t.Errorf("expected CodeInternal, got %s", e.Code)
	}
	if e.Message != "failed to save" {
		t.Errorf("unexpected message: %s", e.Message)
	}
	if e.Cause != cause {
		t.Errorf("expected Cause to be original error, got=%v", e.Cause)
	}
	if e.Details["retryable"] != true {
		t.Errorf("unexpected details: %#v", e.Details)
	}
}

// TestError_Unwrap_ExposesCause は Unwrap で元の cause が取り出せることを確認します。
func TestError_Unwrap_ExposesCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("root")
	e := apperr.Wrap(apperr.CodeInternal, "wrap", cause, nil)

	if !errors.Is(e, cause) {
		t.Errorf("expected errors.Is to match wrapped cause")
	}

	if got := errors.Unwrap(e); got != cause {
		t.Errorf("expected Unwrap()=%v, got=%v", cause, got)
	}
}

//
// ---------- Is ----------
//

// TestIs_MatchCode は 同じ Code の apperr.Error に対して true を返すことを確認します。
func TestIs_MatchCode(t *testing.T) {
	t.Parallel()

	err := apperr.Wrap(apperr.CodeUnauthorized, "unauthorized", errors.New("token expired"), nil)

	if !apperr.Is(err, apperr.CodeUnauthorized) {
		t.Errorf("expected Is(..., CodeUnauthorized)=true")
	}
}

// TestIs_DifferentCode は 異なる Code の場合に false を返すことを確認します。
func TestIs_DifferentCode(t *testing.T) {
	t.Parallel()

	err := apperr.Wrap(apperr.CodeForbidden, "forbidden", errors.New("no permission"), nil)

	if apperr.Is(err, apperr.CodeUnauthorized) {
		t.Errorf("expected Is(..., CodeUnauthorized)=false for forbidden error")
	}
}

// TestIs_NonAppErr は *apperr.Error 以外の error に対して false を返すことを確認します。
func TestIs_NonAppErr(t *testing.T) {
	t.Parallel()

	err := errors.New("plain error")

	if apperr.Is(err, apperr.CodeInternal) {
		t.Errorf("expected Is(..., CodeInternal)=false for non-app error")
	}
}
