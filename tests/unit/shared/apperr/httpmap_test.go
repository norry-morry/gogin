// Package apperr_test は application error code と HTTP ステータスの
// マッピング (internal/shared/apperr/httpmap.go) をテストします。
package apperr_test

import (
	"net/http"
	"testing"

	"resume/internal/shared/apperr"
)

// TestHTTPStatusOf_KnownCodes は 代表的な Code が想定した HTTP ステータスに変換されることを確認します。
func TestHTTPStatusOf_KnownCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code apperr.Code
		want int
	}{
		{apperr.CodeBadRequest, http.StatusBadRequest},
		{apperr.CodeUnprocessable, http.StatusUnprocessableEntity},
		{apperr.CodeUnauthorized, http.StatusUnauthorized},
		{apperr.CodeForbidden, http.StatusForbidden},
		{apperr.CodeNotFound, http.StatusNotFound},
		{apperr.CodeConflict, http.StatusConflict},
		{apperr.CodeInternal, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.code), func(t *testing.T) {
			t.Parallel()

			got := apperr.HTTPStatusOf(tt.code)
			if got != tt.want {
				t.Errorf("HTTPStatusOf(%s) expected %d, got %d", tt.code, tt.want, got)
			}
		})
	}
}

// TestHTTPStatusOf_UnknownCode は 想定外の Code に対して 500 を返すことを確認します。
func TestHTTPStatusOf_UnknownCode(t *testing.T) {
	t.Parallel()

	got := apperr.HTTPStatusOf(apperr.Code("SOMETHING_ELSE"))
	if got != http.StatusInternalServerError {
		t.Errorf("expected 500 for unknown code, got %d", got)
	}
}
