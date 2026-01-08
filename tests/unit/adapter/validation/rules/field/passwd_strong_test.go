// Package field_test は、フィールドレベルのカスタムバリデーション
// （jp_postal / jp_pref / no_edge_spaces / passwd_strong）に対する
// ユニットテストを提供します。
package field_test

import (
	"testing"

	_ "resume/internal/adapter/validation/rules/field" // jp_postal などの init 登録用
)

// passwdStrongPayload はデフォルトの passwd_strong 検証用のテスト構造体です。
// デフォルトでは長さ 12 以上かつ 4 種類のうち 3 種類以上（小文字・大文字・数字・記号）を要求します。
type passwdStrongPayload struct {
	Password string `validate:"passwd_strong"`
}

// passwdStrongPayloadCustom は、最小長をパラメータで指定するテスト用構造体です。
type passwdStrongPayloadCustom struct {
	Password string `validate:"passwd_strong=8"`
}

// TestPasswdStrong_TooShort_NG は、強度要件を満たしていても長さが足りなければエラーになることを確認します。
func TestPasswdStrong_TooShort_NG(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	// 種類は十分だが長さ < 12
	in := passwdStrongPayload{
		Password: "Ab1!xyz", // 7 文字
	}
	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error for too short password, got nil")
	}
}

// TestPasswdStrong_OnlyLower_NG は、長さは足りていても小文字のみではエラーになることを確認します。
func TestPasswdStrong_OnlyLower_NG(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := passwdStrongPayload{
		Password: "abcdefghijklmn", // 14 文字, 1 種類のみ
	}
	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error for lower-only password, got nil")
	}
}

// TestPasswdStrong_Valid_OK は、長さ・文字種ともに条件を満たす場合にエラーにならないことを確認します。
func TestPasswdStrong_Valid_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := passwdStrongPayload{
		Password: "Abcdef1234!@", // 12 文字, lower/upper/digit/symbol の 4 種類
	}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for strong password, got=%v", err)
	}
}

// TestPasswdStrong_CustomMin_OK は、パラメータで最小長を短く指定した場合に
// その長さに応じた検証が行われることを確認します。
func TestPasswdStrong_CustomMin_OK(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := passwdStrongPayloadCustom{
		Password: "Ab1!xyz9", // 8 文字, 種類は 4 種
	}
	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error for strong password with custom min, got=%v", err)
	}
}

// TestPasswdStrong_CustomMin_TooShort_NG は、カスタム最小長を指定しても
// その長さに満たない場合はエラーになることを確認します。
func TestPasswdStrong_CustomMin_TooShort_NG(t *testing.T) {
	t.Parallel()

	v := newValidatorForFieldRules(t)

	in := passwdStrongPayloadCustom{
		Password: "Ab1!xyz", // 7 文字 < 8
	}
	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error for password shorter than custom min, got nil")
	}
}
