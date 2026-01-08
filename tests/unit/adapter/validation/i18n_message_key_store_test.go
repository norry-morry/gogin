// Package validation_test は I18nMessageKeyStore のユニットテストを提供します。
package validation_test

import (
	"strings"
	"testing"

	"resume/internal/adapter/validation"
	vo "resume/internal/domain/valueobject/i18n"
)

// fakeTranslator は uci18n.Translator を満たすテスト用の fake 実装です。
type fakeTranslator struct {
	store map[string]string
}

// Translate は store に登録されたキーがあれば値と true を、
// なければ空文字と false を返します。
func (f *fakeTranslator) Translate(_ vo.Locale, key vo.Key) (string, bool) {
	v, ok := f.store[key.Value()]
	return v, ok
}

// Bundle は prefix でフィルタしたキーのバンドルを返します。
// テストでは HasKey しか使わないため、最低限の実装です。
func (f *fakeTranslator) Bundle(_ vo.Locale, prefix string) map[string]string {
	result := make(map[string]string)
	for k, v := range f.store {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			result[k] = v
		}
	}
	return result
}

// TestI18nMessageKeyStore_HasKey_ReturnsTrueWhenPresent は
// Translator がキーを保持している場合 true が返ることを確認します。
func TestI18nMessageKeyStore_HasKey_ReturnsTrueWhenPresent(t *testing.T) {
	t.Parallel()

	fake := &fakeTranslator{
		store: map[string]string{
			"validation.required": "必須項目です",
		},
	}

	s := validation.NewI18nMessageKeyStore(fake, vo.NewLocale("ja"))

	if !s.HasKey("validation.required") {
		t.Errorf("expected true for existing key, got false")
	}
}

// TestI18nMessageKeyStore_HasKey_ReturnsFalseWhenAbsent は
// Translator がキーを保持していない場合 false が返ることを確認します。
func TestI18nMessageKeyStore_HasKey_ReturnsFalseWhenAbsent(t *testing.T) {
	t.Parallel()

	fake := &fakeTranslator{
		store: map[string]string{}, // 空
	}

	s := validation.NewI18nMessageKeyStore(fake, vo.NewLocale("ja"))

	if s.HasKey("validation.missing") {
		t.Errorf("expected false for missing key, got true")
	}
}
