// Package field_test は、フィールドレベルのカスタムバリデーション
// （jp_postal / jp_pref / no_edge_spaces / passwd_strong）に対する
// ユニットテストを提供します。
package field_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	rules "resume/internal/adapter/validation/rules"
	_ "resume/internal/adapter/validation/rules/field" // jp_postal などの init 登録用
)

// jpPostalPayload は jp_postal 検証用のテスト構造体です。
type jpPostalPayload struct {
	Postal string `validate:"jp_postal"`
}

func newValidatorForFieldRules(t *testing.T) *validator.Validate {
	t.Helper()

	v := validator.New()
	if err := rules.ApplyAll(v); err != nil {
		t.Fatalf("failed to apply validation rules: %v", err)
	}
	return v
}

// TestJPPostal_Empty_OK は、空文字は許可される（omitempty 前提）ことを確認します。
func TestJPPostal_Empty_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := jpPostalPayload{Postal: ""}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for empty postal, got=%v", err)
	}
}

// TestJPPostal_Valid_OK は、NNN-NNNN 形式の郵便番号が許可されることを確認します。
func TestJPPostal_Valid_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := jpPostalPayload{Postal: "123-4567"}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for valid postal, got=%v", err)
	}

	in2 := jpPostalPayload{Postal: "001-0001"}
	if err := v.Struct(in2); err != nil {
		t.Fatalf("expected no error for valid postal with leading zeros, got=%v", err)
	}
}

// TestJPPostal_Invalid_NG は、不正な形式の郵便番号がエラーになることを確認します。
func TestJPPostal_Invalid_NG(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"no hyphen":    "1234567",
		"short":        "12-3456",
		"long":         "1234-5678",
		"alpha":        "abc-defg",
		"wrong format": "12-34567",
	}

	v := newValidatorForFieldRules(t)

	for name, postal := range tests {
		postal := postal
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in := jpPostalPayload{Postal: postal}
			err := v.Struct(in)
			if err == nil {
				t.Fatalf("expected error for postal=%q, got nil", postal)
			}

			verrs, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatalf("expected validator.ValidationErrors, got %T", err)
			}
			if len(verrs) == 0 {
				t.Fatalf("expected at least one validation error")
			}

			if verrs[0].Tag() != "jp_postal" {
				t.Errorf("expected Tag=jp_postal, got=%s", verrs[0].Tag())
			}
		})
	}
}
