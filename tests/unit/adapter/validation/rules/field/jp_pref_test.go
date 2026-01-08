// Package field_test は、フィールドレベルのカスタムバリデーション
// （jp_postal / jp_pref / no_edge_spaces / passwd_strong）に対する
// ユニットテストを提供します。
package field_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	_ "resume/internal/adapter/validation/rules/field" // jp_postal などの init 登録用
)

// jpPrefPayload は jp_pref 検証用のテスト構造体です。
type jpPrefPayload struct {
	Pref string `validate:"jp_pref"`
}

// TestJPPref_Empty_OK は、空文字が許可される（omitempty 前提）ことを確認します。
func TestJPPref_Empty_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := jpPrefPayload{Pref: ""}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for empty prefecture, got=%v", err)
	}
}

// TestJPPref_Valid_OK は、有効な都道府県名が許可されることを確認します。
func TestJPPref_Valid_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	// jpstring.IsValidPrefecture で true になる想定の値
	in := jpPrefPayload{Pref: "東京都"}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for valid prefecture, got=%v", err)
	}
}

// TestJPPref_Invalid_NG は、存在しない都道府県名がエラーになることを確認します。
func TestJPPref_Invalid_NG(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := jpPrefPayload{Pref: "存在しない県名"}
	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error for invalid prefecture, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}
	if len(verrs) == 0 {
		t.Fatalf("expected at least one validation error")
	}

	if verrs[0].Tag() != "jp_pref" {
		t.Errorf("expected Tag=jp_pref, got=%s", verrs[0].Tag())
	}
}
