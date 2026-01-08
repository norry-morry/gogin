package i18n_test

import (
	"testing"

	i18nvo "resume/internal/domain/valueobject/i18n"
)

// TestNewLocale_StoresCode は NewLocale の引数がそのまま保持されることを確認します。
func TestNewLocale_StoresCode(t *testing.T) {
	t.Parallel()

	ja := i18nvo.NewLocale("ja")

	if got := ja.Code(); got != "ja" {
		t.Errorf("expected code=ja, got=%s", got)
	}
}

// TestLocale_Equal_SameCode は 同じコード同士で Equal が true になることを確認します。
func TestLocale_Equal_SameCode(t *testing.T) {
	t.Parallel()

	l1 := i18nvo.NewLocale("en")
	l2 := i18nvo.NewLocale("en")

	if !l1.Equal(l2) {
		t.Errorf("expected locales to be equal")
	}
}

// TestLocale_Equal_DifferentCode は 異なるコード同士で Equal が false になることを確認します。
func TestLocale_Equal_DifferentCode(t *testing.T) {
	t.Parallel()

	l1 := i18nvo.NewLocale("ja")
	l2 := i18nvo.NewLocale("en")

	if l1.Equal(l2) {
		t.Errorf("expected locales to be not equal")
	}
}

// TestLocale_String_EqualsCode は String() がコードと同じ文字列を返すことを確認します。
func TestLocale_String_EqualsCode(t *testing.T) {
	t.Parallel()

	l := i18nvo.NewLocale("fr")

	if got := l.String(); got != "fr" {
		t.Errorf("expected String()=fr, got=%s", got)
	}
}
