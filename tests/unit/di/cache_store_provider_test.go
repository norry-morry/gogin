// Package di_test は internal/di パッケージで定義された
// DI プロバイダ群のうち、i18n のキャッシュストア関連 provider の
// 挙動を検証するユニットテストを提供します。
package di_test

import (
	"testing"

	gatewaymem "resume/internal/adapter/gateway/inmemory"
	"resume/internal/di"
	uci18n "resume/internal/usecase/i18n"
)

// TestExportProvideCacheStoreImpl_ReturnsConcrete は、
// ExportProvideCacheStoreImpl が *gatewaymem.CacheStoreImpl を返すことを確認します。
func TestExportProvideCacheStoreImpl_ReturnsConcrete(t *testing.T) {
	t.Parallel()

	cs := di.ExportProvideCacheStoreImpl()
	if cs == nil {
		t.Fatalf("expected non-nil CacheStoreImpl from ExportProvideCacheStoreImpl")
	}

	// 具体型として *gatewaymem.CacheStoreImpl であることを確認
	if _, ok := any(cs).(*gatewaymem.CacheStoreImpl); !ok {
		t.Fatalf("expected *gatewaymem.CacheStoreImpl, got %T", cs)
	}
}

// TestExportProvideTranslator_UsesSameInstance は、
// ExportProvideTranslator が CacheStoreImpl と同一インスタンスを
// uci18n.Translator として返すことを確認します。
func TestExportProvideTranslator_UsesSameInstance(t *testing.T) {
	t.Parallel()

	// まず具象の CacheStoreImpl を取得
	cs := di.ExportProvideCacheStoreImpl()

	tr := di.ExportProvideTranslator(cs)
	if tr == nil {
		t.Fatalf("expected non-nil Translator from ExportProvideTranslator")
	}

	// インターフェース満たしていることをコンパイル時に保証
	var _ uci18n.Translator = tr

	// 実際に同一ポインタであることを確認
	if tr != cs {
		t.Errorf("expected translator and cacheStore to be same instance, got different pointers")
	}
}

// TestExportProvideCacheStore_UsesSameInstance は、
// ExportProvideCacheStore が CacheStoreImpl と同一インスタンスを
// uci18n.CacheStore として返すことを確認します。
func TestExportProvideCacheStore_UsesSameInstance(t *testing.T) {
	t.Parallel()

	cs := di.ExportProvideCacheStoreImpl()

	store := di.ExportProvideCacheStore(cs)
	if store == nil {
		t.Fatalf("expected non-nil CacheStore from ExportProvideCacheStore")
	}

	// インターフェース満たしていることをコンパイル時に保証
	var _ uci18n.CacheStore = store

	// 実際に同一ポインタであることを確認
	if store != cs {
		t.Errorf("expected cacheStore and original CacheStoreImpl to be same instance, got different pointers")
	}
}
