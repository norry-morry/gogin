// Package field_test は、フィールドレベルのカスタムバリデーション
// （jp_postal / jp_pref / no_edge_spaces / passwd_strong）に対する
// ユニットテストを提供します。
package field_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	_ "resume/internal/adapter/validation/rules/field" // jp_postal などの init 登録用
)

// noEdgeSpacesPayload は no_edge_spaces 検証用のテスト構造体です。
type noEdgeSpacesPayload struct {
	Val string `validate:"no_edge_spaces"`
}

// TestNoEdgeSpaces_Valid_OK は、前後に空白がない場合にエラーにならないことを確認します。
func TestNoEdgeSpaces_Valid_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	cases := []string{
		"",
		"abc",
		"a b", // ← OK（前後スペースなし）
	}

	for _, s := range cases {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()

			in := noEdgeSpacesPayload{Val: s}
			if err := v.Struct(in); err != nil {
				t.Fatalf("expected no error for %q, got=%v", s, err)
			}
		})
	}
}

// TestNoEdgeSpaces_Invalid_NG は、前後に空白がある場合にエラーとなることを確認します。
func TestNoEdgeSpaces_Invalid_NG(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	cases := []string{
		" abc",
		"abc ",
		"  abc",
		"abc  ",
		" abc ",
	}

	for _, s := range cases {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()

			in := noEdgeSpacesPayload{Val: s}
			err := v.Struct(in)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", s)
			}

			verrs, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatalf("expected validator.ValidationErrors, got %T", err)
			}
			if len(verrs) == 0 {
				t.Fatalf("expected at least one validation error")
			}

			if verrs[0].Tag() != "no_edge_spaces" {
				t.Errorf("expected Tag=no_edge_spaces, got=%s", verrs[0].Tag())
			}
		})
	}
}
