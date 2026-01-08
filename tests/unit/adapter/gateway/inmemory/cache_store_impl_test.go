// Package inmemory_test は internal/adapter/gateway/inmemory パッケージに
// 定義された i18n キャッシュストアのユニットテストを提供します。
package inmemory_test

import (
	"testing"

	inmemory "resume/internal/adapter/gateway/inmemory"
	vo "resume/internal/domain/valueobject/i18n"
)

// helper: Locale / Key / Bundle を簡単に生成する
func loc(s string) vo.Locale { return vo.NewLocale(s) }
func key(s string) vo.Key    { return vo.NewKey(s) }

func bundle(locale string, msgs map[string]string) *vo.Bundle {
	return vo.NewBundle(loc(locale), msgs)
}

// -----------------------------------------------------------------------------
// NewCacheStoreImpl
// -----------------------------------------------------------------------------
func TestNewCacheStoreImpl_InitEmptySnapshot(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	locs := cs.Locales()

	if len(locs) != 0 {
		t.Fatalf("expected no locales initially, got=%v", locs)
	}
}

// -----------------------------------------------------------------------------
// SwapAll
// -----------------------------------------------------------------------------
func TestCacheStore_SwapAll_ReplacesSnapshot(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()

	all := map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{"a": "あ"}),
		loc("en"): bundle("en", map[string]string{"a": "A"}),
	}

	if err := cs.SwapAll(all); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	got, ok := cs.Get(loc("ja"), key("a"))
	if !ok || got != "あ" {
		t.Fatalf("expected ja:a=あ, got=%v ok=%v", got, ok)
	}

	got, ok = cs.Get(loc("en"), key("a"))
	if !ok || got != "A" {
		t.Fatalf("expected en:a=A, got=%v ok=%v", got, ok)
	}
}

// -----------------------------------------------------------------------------
// Get
// -----------------------------------------------------------------------------
func TestCacheStore_Get_NotFound(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	val, ok := cs.Get(loc("ja"), key("x"))
	if ok || val != "" {
		t.Fatalf("expected not found, got=%v ok=%v", val, ok)
	}
}

// -----------------------------------------------------------------------------
// Set
// -----------------------------------------------------------------------------
func TestCacheStore_Set_UpdatesExistingLocale(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{"a": "initial"}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	cs.Set(loc("ja"), key("a"), "updated")

	got, ok := cs.Get(loc("ja"), key("a"))
	if !ok || got != "updated" {
		t.Fatalf("expected updated value, got=%v ok=%v", got, ok)
	}
}

func TestCacheStore_Set_CreatesLocaleIfMissing(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()

	cs.Set(loc("en"), key("hello"), "world")

	got, ok := cs.Get(loc("en"), key("hello"))
	if !ok || got != "world" {
		t.Fatalf("expected new locale value, got=%v ok=%v", got, ok)
	}
}

// -----------------------------------------------------------------------------
// Translate（フォールバック）
// -----------------------------------------------------------------------------
func TestCacheStore_Translate_PrimaryLocaleFound(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{"hello": "こんにちは"}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	got, ok := cs.Translate(loc("ja"), key("hello"))
	if !ok || got != "こんにちは" {
		t.Fatalf("expected ja translation, got=%v ok=%v", got, ok)
	}
}

func TestCacheStore_Translate_FallbackToEnglish(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("en"): bundle("en", map[string]string{"hello": "Hello"}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	// ja にはないが en にはあるのでフォールバック
	got, ok := cs.Translate(loc("ja"), key("hello"))
	if !ok || got != "Hello" {
		t.Fatalf("expected English fallback, got=%v ok=%v", got, ok)
	}
}

func TestCacheStore_Translate_NotFoundAnywhere(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()

	got, ok := cs.Translate(loc("ja"), key("none"))
	if ok || got != "" {
		t.Fatalf("expected not found, got=%v ok=%v", got, ok)
	}
}

// -----------------------------------------------------------------------------
// Bundle（prefix）
// -----------------------------------------------------------------------------
func TestCacheStore_Bundle_PrefixMatch(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{
			"a.b.c": "1",
			"a.b.d": "2",
			"x":     "3",
		}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	res := cs.Bundle(loc("ja"), "a.b.")

	if len(res) != 2 {
		t.Fatalf("expected 2 matches, got=%v", res)
	}
	if res["a.b.c"] != "1" || res["a.b.d"] != "2" {
		t.Fatalf("unexpected bundle: %v", res)
	}
}

func TestCacheStore_Bundle_EmptyPrefix_ReturnsAll(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{
			"a": "1",
			"b": "2",
		}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	res := cs.Bundle(loc("ja"), "")

	if len(res) != 2 {
		t.Fatalf("expected 2 messages, got=%v", res)
	}
}

// -----------------------------------------------------------------------------
// Locales / Clear
// -----------------------------------------------------------------------------
func TestCacheStore_Locales(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{"x": "1"}),
		loc("en"): bundle("en", map[string]string{"x": "1"}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	locs := cs.Locales()
	if len(locs) != 2 {
		t.Fatalf("expected 2 locales, got=%v", locs)
	}
}

func TestCacheStore_Clear(t *testing.T) {
	t.Parallel()

	cs := inmemory.NewCacheStoreImpl()
	if err := cs.SwapAll(map[vo.Locale]*vo.Bundle{
		loc("ja"): bundle("ja", map[string]string{"a": "1"}),
	}); err != nil {
		t.Fatalf("SwapAll returned error: %v", err)
	}

	cs.Clear()

	if locs := cs.Locales(); len(locs) != 0 {
		t.Fatalf("expected empty locales after Clear, got=%v", locs)
	}
}
