// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルでは i18n 用 DictionaryRepository
// プロバイダのテストを行います。
package di_test

import (
	"testing"

	"resume/internal/di"
	repo "resume/internal/domain/repository"
)

// TestProvideDictionaryRepository_ReturnsImplementation は ExportProvideDictionaryRepository が
// DictionaryRepository 実装を返すことを確認します。
func TestProvideDictionaryRepository_ReturnsImplementation(t *testing.T) {
	t.Parallel()

	r := di.ExportProvideDictionaryRepository("")

	if r == nil {
		t.Fatalf("expected non-nil DictionaryRepository from ExportProvideDictionaryRepository")
	}

	// DictionaryRepository インタフェースとして扱えることをコンパイル時に保証
	var _ repo.DictionaryRepository = r
}
