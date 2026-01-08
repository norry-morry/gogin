// Package jp_test は、日本ロケール向けのバリデーションラッパ
// WrapWhenJP に対するユニットテストを提供します。
package jp_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	jprules "resume/internal/adapter/validation/rules/locale/jp"
)

// addressPayloadJP は WrapWhenJP のテスト用構造体です。
// CountryCode および AdministrativeArea フィールド名は
// 本実装の WrapWhenJP が reflect で参照する名前に合わせています。
type addressPayloadJP struct {
	CountryCode        string
	AdministrativeArea string
}

// newValidatorWithWhenJP は、WrapWhenJP を struct-level バリデーションとして
// 登録した validator.Validate を返します。
func newValidatorWithWhenJP(tag string) *validator.Validate {
	v := validator.New()
	v.RegisterStructValidation(jprules.WrapWhenJP(tag), addressPayloadJP{})
	return v
}

// TestWrapWhenJP_NotJP_NoError は、CountryCode が JP 以外の場合は
// AdministrativeArea が空でもエラーが発生しないことを確認します。
func TestWrapWhenJP_NotJP_NoError(t *testing.T) {
	t.Parallel()

	v := newValidatorWithWhenJP("jp_pref_required")

	in := addressPayloadJP{
		CountryCode:        "US",
		AdministrativeArea: "",
	}

	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error when CountryCode is not JP, got=%v", err)
	}
}

// TestWrapWhenJP_JP_WithPref_NoError は、CountryCode=JP かつ
// AdministrativeArea が非空の場合はエラーにならないことを確認します。
func TestWrapWhenJP_JP_WithPref_NoError(t *testing.T) {
	t.Parallel()

	v := newValidatorWithWhenJP("jp_pref_required")

	in := addressPayloadJP{
		CountryCode:        "JP",
		AdministrativeArea: "東京都",
	}

	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error when JP and AdministrativeArea is set, got=%v", err)
	}
}

// TestWrapWhenJP_JP_WithoutPref_Error は、CountryCode=JP かつ
// AdministrativeArea が空の場合に、指定したタグでエラーが出ることを確認します。
func TestWrapWhenJP_JP_WithoutPref_Error(t *testing.T) {
	t.Parallel()

	const tag = "jp_pref_required"

	v := newValidatorWithWhenJP(tag)

	in := addressPayloadJP{
		CountryCode:        "JP",
		AdministrativeArea: "",
	}

	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error when JP and AdministrativeArea is empty, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}
	if len(verrs) == 0 {
		t.Fatalf("expected at least 1 validation error, got 0")
	}

	// AdministrativeArea に対するエラーがあり、タグと Param が期待通りか確認する
	found := false
	for _, e := range verrs {
		if e.Tag() == tag && e.Param() == "required" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected tag=%s error, but none found", tag)
	}
}

// TestWrapWhenJP_CaseInsensitiveCountryCode は、CountryCode が "jp" のような
// 小文字でも JP と同等に扱われることを確認します。
func TestWrapWhenJP_CaseInsensitiveCountryCode(t *testing.T) {
	t.Parallel()

	const tag = "jp_pref_required"

	v := newValidatorWithWhenJP(tag)

	in := addressPayloadJP{
		CountryCode:        "jp", // 小文字
		AdministrativeArea: "",
	}

	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error when CountryCode is 'jp' (lowercase) and AdministrativeArea is empty, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	found := false
	for _, e := range verrs {
		if e.Tag() == tag && e.Param() == "required" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected tag=%s error, but none found", tag)
	}
}
