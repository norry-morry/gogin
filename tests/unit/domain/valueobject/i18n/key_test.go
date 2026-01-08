// Package i18n_test は i18n 値オブジェクト (Key, Locale, Bundle) の
// ユニットテストを提供します。
package i18n_test

import (
	"testing"

	i18nvo "resume/internal/domain/valueobject/i18n"
)

// TestNewKey_StoresValue は NewKey の引数がそのまま保持されることを確認します。
func TestNewKey_StoresValue(t *testing.T) {
	t.Parallel()

	k := i18nvo.NewKey("ui.page.profile.title")

	if got := k.Value(); got != "ui.page.profile.title" {
		t.Errorf("expected value=ui.page.profile.title, got=%s", got)
	}
}

// TestKey_Equal_SameValue は 同じ値の Key 同士が Equal で true を返すことを確認します。
func TestKey_Equal_SameValue(t *testing.T) {
	t.Parallel()

	k1 := i18nvo.NewKey("master.country.jp")
	k2 := i18nvo.NewKey("master.country.jp")

	if !k1.Equal(k2) {
		t.Errorf("expected keys to be equal")
	}
}

// TestKey_Equal_DifferentValue は 異なる値の Key 同士が Equal で false を返すことを確認します。
func TestKey_Equal_DifferentValue(t *testing.T) {
	t.Parallel()

	k1 := i18nvo.NewKey("master.country.jp")
	k2 := i18nvo.NewKey("master.country.us")

	if k1.Equal(k2) {
		t.Errorf("expected keys to be not equal")
	}
}

// TestKey_String_EqualsValue は String() が内部値と同じ文字列を返すことを確認します。
func TestKey_String_EqualsValue(t *testing.T) {
	t.Parallel()

	k := i18nvo.NewKey("ui.page.dashboard.title")

	if got := k.String(); got != "ui.page.dashboard.title" {
		t.Errorf("expected String()=ui.page.dashboard.title, got=%s", got)
	}
}
