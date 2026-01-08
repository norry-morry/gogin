package i18n_test

import (
	"testing"

	i18nvo "resume/internal/domain/valueobject/i18n"
)

// TestNewBundle_StoresLocaleAndMessages は Locale とメッセージが格納されることを確認します。
func TestNewBundle_StoresLocaleAndMessages(t *testing.T) {
	t.Parallel()

	locale := i18nvo.NewLocale("ja")
	src := map[string]string{
		"ui.page.dashboard.title": "ダッシュボード",
		"ui.page.profile.title":   "プロフィール",
	}

	b := i18nvo.NewBundle(locale, src)

	// Locale がそのまま保持されているか
	if !b.Locale().Equal(locale) {
		t.Errorf("expected bundle locale=%s, got=%s", locale.Code(), b.Locale().Code())
	}

	// Get で取得できるか
	title, ok := b.Get(i18nvo.NewKey("ui.page.dashboard.title"))
	if !ok {
		t.Fatalf("expected key ui.page.dashboard.title to exist")
	}
	if title != "ダッシュボード" {
		t.Errorf("expected ダッシュボード, got=%s", title)
	}
}

// TestBundle_Get_NotFound は 存在しないキーでは ok=false になることを確認します。
func TestBundle_Get_NotFound(t *testing.T) {
	t.Parallel()

	locale := i18nvo.NewLocale("ja")
	src := map[string]string{
		"ui.page.dashboard.title": "ダッシュボード",
	}
	b := i18nvo.NewBundle(locale, src)

	_, ok := b.Get(i18nvo.NewKey("ui.page.profile.title"))
	if ok {
		t.Errorf("expected ok=false for missing key")
	}
}

// TestBundle_Immutability_InputMap は 元の map を変更しても Bundle 側に影響しないことを確認します。
func TestBundle_Immutability_InputMap(t *testing.T) {
	t.Parallel()

	locale := i18nvo.NewLocale("ja")
	src := map[string]string{
		"ui.page.dashboard.title": "ダッシュボード",
	}
	b := i18nvo.NewBundle(locale, src)

	// NewBundle 生成後に元の map を書き換え
	src["ui.page.dashboard.title"] = "書き換え後"

	v, _ := b.Get(i18nvo.NewKey("ui.page.dashboard.title"))
	if v != "ダッシュボード" {
		t.Errorf("expected bundle to keep original value, got=%s", v)
	}
}

// TestBundle_Messages_ReturnsCopy は Messages() から返された map を変更しても
// 内部状態に影響しないことを確認します。
func TestBundle_Messages_ReturnsCopy(t *testing.T) {
	t.Parallel()

	locale := i18nvo.NewLocale("ja")
	src := map[string]string{
		"ui.page.dashboard.title": "ダッシュボード",
	}
	b := i18nvo.NewBundle(locale, src)

	msgs := b.Messages()
	msgs["ui.page.dashboard.title"] = "書き換え後（外側）"

	// Bundle 内の値は変わらないはず
	v, _ := b.Get(i18nvo.NewKey("ui.page.dashboard.title"))
	if v != "ダッシュボード" {
		t.Errorf("expected bundle to be immutable from Messages(), got=%s", v)
	}
}
