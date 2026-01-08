// Package rules_test は internal/adapter/validation/rules パッケージに定義された
// バリデーションエラー → i18n 対応エラーマップ変換処理のユニットテストを提供します。
package rules_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
	"resume/internal/shared/util"
)

// fakeMessageKeyStore は rules.MessageKeyStore を満たすテスト用の fake 実装です。
type fakeMessageKeyStore struct {
	keys map[string]bool
}

func (f *fakeMessageKeyStore) HasKey(key string) bool {
	return f.keys[key]
}

// helper: ValidationErrors を作るための簡易 struct
type requiredPayload struct {
	Name string `validate:"required"`
}

type minPayload struct {
	Age int `validate:"min=3"`
}

// TestMapValidationErrors_UsesRulesKeyWhenNoStore は、
// MessageKeyStore が設定されていない場合に rules キーが使われることを確認します。
func TestMapValidationErrors_UsesRulesKeyWhenNoStore(t *testing.T) {
	v := validator.New()

	// Name が未設定なので "required" エラーが発生する
	var p requiredPayload
	err := v.Struct(p)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	// MessageKeyStore を明示的に nil にする
	rules.SetMessageKeyStore(nil)

	scope := "domain.user"
	out := rules.MapValidationErrors(verrs, scope)

	fieldLowerCamel := util.ToCamel("Name")
	base := scope + "." + fieldLowerCamel

	items, ok := out[base]
	if !ok || len(items) == 0 {
		t.Fatalf("expected error items for key %s", base)
	}

	got := items[0]
	wantCode := "validation.rules.required"

	if got.Code != wantCode {
		t.Errorf("expected Code=%s, got=%s", wantCode, got.Code)
	}

	// field パラメータが labelKey になっていることを確認
	labelKey := base + ".label"
	if got.Params["field"] != labelKey {
		t.Errorf("expected field param=%s, got=%v", labelKey, got.Params["field"])
	}

	// "required" は Param を持たないので "param" キーは存在しないはず
	if _, exists := got.Params["param"]; exists {
		t.Errorf("did not expect param in Params for required rule")
	}
}

// TestMapValidationErrors_PrefersCustomKeyOverRules は、
// custom 用のキーが存在する場合に rules より custom キーが優先されることを確認します。
func TestMapValidationErrors_PrefersCustomKeyOverRules(t *testing.T) {
	v := validator.New()

	var p requiredPayload
	err := v.Struct(p)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	// custom 用のキーだけを持つ MessageKeyStore をセット
	store := &fakeMessageKeyStore{
		keys: map[string]bool{
			"validation.custom.required": true,
		},
	}
	rules.SetMessageKeyStore(store)

	scope := "domain.user"
	out := rules.MapValidationErrors(verrs, scope)

	fieldLowerCamel := util.ToCamel("Name")
	base := scope + "." + fieldLowerCamel

	items, ok := out[base]
	if !ok || len(items) == 0 {
		t.Fatalf("expected error items for key %s", base)
	}

	got := items[0]
	wantCode := "validation.custom.required"

	if got.Code != wantCode {
		t.Errorf("expected Code=%s, got=%s", wantCode, got.Code)
	}
}

// TestMapValidationErrors_PrefersFieldSpecificOverCustom は、
// フィールド固有のキーが存在する場合に custom キーより優先されることを確認します。
func TestMapValidationErrors_PrefersFieldSpecificOverCustom(t *testing.T) {
	v := validator.New()

	var p requiredPayload
	err := v.Struct(p)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	scope := "domain.user"
	fieldLowerCamel := util.ToCamel("Name")

	// lookup 用のキーは snake ベース:
	// fieldSpecificKey(scope, fieldLowerCamel, "required") の "required" 部分が snake
	fsLookupKey := "validation.fieldspecific." + scope + "." + fieldLowerCamel + ".required"

	store := &fakeMessageKeyStore{
		keys: map[string]bool{
			fsLookupKey:                  true,
			"validation.custom.required": true, // あっても fieldSpecific が優先されるはず
		},
	}
	rules.SetMessageKeyStore(store)

	out := rules.MapValidationErrors(verrs, scope)

	base := scope + "." + fieldLowerCamel
	items, ok := out[base]
	if !ok || len(items) == 0 {
		t.Fatalf("expected error items for key %s", base)
	}

	got := items[0]
	// 表示用キーは lowerCamel を使う:
	// fieldSpecificKey(scope, fieldLowerCamel, "required")
	wantCode := "validation.fieldspecific." + scope + "." + fieldLowerCamel + ".required"

	if got.Code != wantCode {
		t.Errorf("expected Code=%s, got=%s", wantCode, got.Code)
	}
}

// TestMapValidationErrors_SetsParamsForMinRule は、min ルールの Param/Min が
// 正しく埋め込まれることを確認します。
func TestMapValidationErrors_SetsParamsForMinRule(t *testing.T) {
	v := validator.New()

	p := minPayload{
		Age: 1, // min=3 に違反
	}
	err := v.Struct(p)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	// MessageKeyStore は nil（rules キーが使われる）
	rules.SetMessageKeyStore(nil)

	scope := "domain.user"
	out := rules.MapValidationErrors(verrs, scope)

	fieldLowerCamel := util.ToCamel("Age")
	base := scope + "." + fieldLowerCamel

	items, ok := out[base]
	if !ok || len(items) == 0 {
		t.Fatalf("expected error items for key %s", base)
	}

	got := items[0]
	if got.Code != "validation.rules.min" {
		t.Errorf("expected Code=validation.rules.min, got=%s", got.Code)
	}

	// param と min の両方に 3 が入っていること
	if got.Params["param"] != "3" {
		t.Errorf("expected param=3, got=%v", got.Params["param"])
	}
	if got.Params["min"] != "3" {
		t.Errorf("expected min=3, got=%v", got.Params["min"])
	}
}
