// Package util_test は機微情報マスキングユーティリティ
// (internal/shared/util/sensitive.go) のユニットテストを提供する。
package util_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"resume/internal/shared/util"
)

//
// ---------- LoadSensitiveConfig ----------
//

// TestLoadSensitiveConfig_Success は YAML ファイルから正常に設定を読み込めることを確認する。
func TestLoadSensitiveConfig_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sensitive.yaml")

	yml := `
sensitive:
  full:
    - password
    - api_key
  part:
    - email
partial:
  keep_head: 2
  keep_tail: 3
  mask_char: "#"
`
	if err := os.WriteFile(path, []byte(yml), 0o644); err != nil {
		t.Fatalf("failed to write temp yaml: %v", err)
	}

	cfg, err := util.LoadSensitiveConfig(path)
	if err != nil {
		t.Fatalf("LoadSensitiveConfig returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if !reflect.DeepEqual(cfg.Sensitive.Full, []string{"password", "api_key"}) {
		t.Errorf("unexpected Sensitive.Full: %#v", cfg.Sensitive.Full)
	}
	if !reflect.DeepEqual(cfg.Sensitive.Part, []string{"email"}) {
		t.Errorf("unexpected Sensitive.Part: %#v", cfg.Sensitive.Part)
	}
	if cfg.Partial.KeepHead != 2 || cfg.Partial.KeepTail != 3 || cfg.Partial.MaskChar != "#" {
		t.Errorf("unexpected Partial config: %#v", cfg.Partial)
	}
}

// TestLoadSensitiveConfig_NotFound は存在しないパスを指定した場合にエラーになることを確認する。
func TestLoadSensitiveConfig_NotFound(t *testing.T) {
	t.Parallel()

	cfg, err := util.LoadSensitiveConfig("no_such_file.yaml")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("expected nil config on error, got %#v", cfg)
	}
}

// TestLoadSensitiveConfig_Defaults は KeepHead/KeepTail が負値、MaskChar が空の場合にデフォルト値が補完されることを確認する。
func TestLoadSensitiveConfig_Defaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sensitive_defaults.yaml")

	yml := `
sensitive:
  full: []
  part: []
partial:
  keep_head: -1
  keep_tail: -1
  mask_char: ""
`
	if err := os.WriteFile(path, []byte(yml), 0o644); err != nil {
		t.Fatalf("failed to write temp yaml: %v", err)
	}

	cfg, err := util.LoadSensitiveConfig(path)
	if err != nil {
		t.Fatalf("LoadSensitiveConfig returned error: %v", err)
	}
	if cfg.Partial.KeepHead != 1 {
		t.Errorf("expected KeepHead=1, got=%d", cfg.Partial.KeepHead)
	}
	if cfg.Partial.KeepTail != 4 {
		t.Errorf("expected KeepTail=4, got=%d", cfg.Partial.KeepTail)
	}
	if cfg.Partial.MaskChar != "*" {
		t.Errorf("expected MaskChar=\"*\", got=%q", cfg.Partial.MaskChar)
	}
}

//
// ---------- FullMask / PartialMask ----------
//

// TestFullMask_AlwaysRedacted は FullMask が常に固定文字列 [REDACTED] を返すことを確認する。
func TestFullMask_AlwaysRedacted(t *testing.T) {
	t.Parallel()

	if got := util.FullMask("anything"); got != util.Redacted {
		t.Errorf("expected %q, got=%q", util.Redacted, got)
	}
}

// TestPartialMask_Basic は部分マスクの基本動作を確認する。
func TestPartialMask_Basic(t *testing.T) {
	t.Parallel()

	s := "123456789"
	got := util.PartialMask(s, 2, 3, "#")
	want := "12####789" // 2 + 4(#) + 3

	if got != want {
		t.Errorf("expected %q, got=%q", want, got)
	}
}

// TestPartialMask_EmptyString は空文字の場合にそのまま返すことを確認する。
func TestPartialMask_EmptyString(t *testing.T) {
	t.Parallel()

	if got := util.PartialMask("", 2, 3, "#"); got != "" {
		t.Errorf("expected empty string, got=%q", got)
	}
}

// TestPartialMask_DefaultMaskChar は maskChar が空の場合に "*" が使用されることを確認する。
func TestPartialMask_DefaultMaskChar(t *testing.T) {
	t.Parallel()

	s := "abcdefg"
	got := util.PartialMask(s, 1, 2, "")
	want := "a****fg" // 1 + 3(*) + 2

	if got != want {
		t.Errorf("expected %q, got=%q", want, got)
	}
}

// TestPartialMask_NegativeKeep は負の keepHead/keepTail が 0 に補正されることを確認する。
func TestPartialMask_NegativeKeep(t *testing.T) {
	t.Parallel()

	s := "secret"
	got := util.PartialMask(s, -1, 2, "*")
	// keepHead=0, keepTail=2 → ****et
	want := "****et"

	if got != want {
		t.Errorf("expected %q, got=%q", want, got)
	}
}

// TestPartialMask_TooShort は keepHead+keepTail >= len(s) の場合にそのまま返すことを確認する。
func TestPartialMask_TooShort(t *testing.T) {
	t.Parallel()

	s := "abc"
	got := util.PartialMask(s, 2, 2, "*")

	if got != s {
		t.Errorf("expected %q, got=%q", s, got)
	}
}

//
// ---------- NewMasker / Masker 実装 (keyMatcher 経由) ----------
//

// TestNewMasker_NilConfigError は nil を渡した場合にエラーとなることを確認する。
func TestNewMasker_NilConfigError(t *testing.T) {
	t.Parallel()

	m, err := util.NewMasker(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// 実装によってはエラーと同時にデフォルト Masker を返す場合もあるため、
	// ここでは m が nil かどうかまでは厳密にチェックしない。
	_ = m
}

// TestMasker_MaskByKey は Full/Part/対象外キーに対するマスク動作を確認する。
func TestMasker_MaskByKey(t *testing.T) {
	t.Parallel()

	cfg := &util.SensitiveConfig{}
	cfg.Sensitive.Full = []string{"password", "token"}
	cfg.Sensitive.Part = []string{"email"}
	cfg.Partial.KeepHead = 2
	cfg.Partial.KeepTail = 3
	cfg.Partial.MaskChar = "*"

	m, err := util.NewMasker(cfg)
	if err != nil {
		t.Fatalf("NewMasker returned error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Masker")
	}

	t.Run("full masked key", func(t *testing.T) {
		got := m.MaskByKey("password", "secret")
		if got != util.Redacted {
			t.Errorf("expected %q, got=%v", util.Redacted, got)
		}
	})

	t.Run("partial masked key", func(t *testing.T) {
		got := m.MaskByKey("email", "user@example.com")
		// 部分マスクされていればよい（完全一致までは要求しない）
		if got == "user@example.com" {
			t.Errorf("expected masked value, got original: %v", got)
		}
	})

	t.Run("non target key", func(t *testing.T) {
		val := "visible"
		got := m.MaskByKey("username", val)
		if got != val {
			t.Errorf("expected value to be unchanged, got=%v", got)
		}
	})
}

//
// ---------- MaskMapShallow ----------
//

// TestMasker_MaskMapShallow は 1 階層の map に対してキーごとのマスクが適用されることを確認する。
func TestMasker_MaskMapShallow(t *testing.T) {
	t.Parallel()

	cfg := &util.SensitiveConfig{}
	cfg.Sensitive.Full = []string{"password"}
	cfg.Sensitive.Part = []string{"email"}
	cfg.Partial.KeepHead = 1
	cfg.Partial.KeepTail = 2
	cfg.Partial.MaskChar = "*"

	m, err := util.NewMasker(cfg)
	if err != nil {
		t.Fatalf("NewMasker returned error: %v", err)
	}

	src := map[string]any{
		"password": "secret",
		"email":    "user@example.com",
		"name":     "alice",
	}

	dst := m.MaskMapShallow(src)

	// 元の map が書き換えられていないこと
	if !reflect.DeepEqual(src["password"], "secret") {
		t.Errorf("expected src[password] to remain secret, got=%v", src["password"])
	}

	// マスク後の map の検証
	if dst["password"] != util.Redacted {
		t.Errorf("expected password to be redacted, got=%v", dst["password"])
	}
	if dst["email"] == src["email"] {
		t.Errorf("expected email to be masked, got=%v", dst["email"])
	}
	if dst["name"] != "alice" {
		t.Errorf("expected name to be unchanged, got=%v", dst["name"])
	}
}

// TestMasker_MaskMapShallow_Nil は nil map を渡した場合に nil が返ることを確認する。
func TestMasker_MaskMapShallow_Nil(t *testing.T) {
	t.Parallel()

	cfg := &util.SensitiveConfig{}
	cfg.Partial.KeepHead = 1
	cfg.Partial.KeepTail = 2
	cfg.Partial.MaskChar = "*"

	m, err := util.NewMasker(cfg)
	if err != nil {
		t.Fatalf("NewMasker returned error: %v", err)
	}

	if got := m.MaskMapShallow(nil); got != nil {
		t.Errorf("expected nil, got=%v", got)
	}
}

//
// ---------- MaskAnyRecursive ----------
//

// TestMasker_MaskAnyRecursive はネストした構造（map / slice）に対して再帰的にマスクが適用されることを確認する。
func TestMasker_MaskAnyRecursive(t *testing.T) {
	t.Parallel()

	cfg := &util.SensitiveConfig{}
	cfg.Sensitive.Full = []string{"password", "token"}
	cfg.Sensitive.Part = []string{"email"}
	cfg.Partial.KeepHead = 1
	cfg.Partial.KeepTail = 2
	cfg.Partial.MaskChar = "*"

	m, err := util.NewMasker(cfg)
	if err != nil {
		t.Fatalf("NewMasker returned error: %v", err)
	}

	orig := map[string]any{
		"password": "secret",
		"profile": map[string]any{
			"email": "user@example.com",
			"nested": []any{
				map[string]any{
					"token": "abc123",
				},
			},
		},
	}

	maskedAny := m.MaskAnyRecursive("", orig)
	masked, ok := maskedAny.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", maskedAny)
	}

	// password は完全マスク
	if masked["password"] != util.Redacted {
		t.Errorf("expected password to be redacted, got=%v", masked["password"])
	}

	// email は部分マスク（値が変わっていればよい）
	profile, ok := masked["profile"].(map[string]any)
	if !ok {
		t.Fatalf("expected profile to be map[string]any, got %T", masked["profile"])
	}

	if profile["email"] == "user@example.com" {
		t.Errorf("expected email to be masked, got=%v", profile["email"])
	}

	// nested の token も完全マスク
	nestedSlice, ok := profile["nested"].([]any)
	if !ok || len(nestedSlice) != 1 {
		t.Fatalf("expected nested to be []any with len=1, got=%#v", profile["nested"])
	}
	nestedMap, ok := nestedSlice[0].(map[string]any)
	if !ok {
		t.Fatalf("expected nested[0] to be map[string]any, got %T", nestedSlice[0])
	}
	if nestedMap["token"] != util.Redacted {
		t.Errorf("expected token to be redacted, got=%v", nestedMap["token"])
	}

	// 元データが書き換えられていないことも確認
	if orig["password"] != "secret" {
		t.Errorf("expected original password to remain secret, got=%v", orig["password"])
	}
	profileVal, ok := orig["profile"].(map[string]any)
	if !ok {
		t.Fatalf("orig[\"profile\"] is not map[string]any: %#v", orig["profile"])
	}
	profileOrig := profileVal
	if profileOrig["email"] != "user@example.com" {
		t.Errorf("expected original email to remain unmasked, got=%v", profileOrig["email"])
	}
}
