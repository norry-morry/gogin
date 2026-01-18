// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルではトランザクションランナー用の
// プロバイダのテストを行います。
package di_test

import (
	"testing"

	"gorm.io/gorm"

	"resume/internal/di"
	stx "resume/internal/shared/tx"
)

// TestProvideSharedTxRunner_ReturnsRunner は ExportProvideSharedTxRunner が
// stx.Runner を返すことを確認します。
func TestProvideSharedTxRunner_ReturnsRunner(t *testing.T) {
	t.Parallel()

	db := &gorm.DB{} // 実際の接続は不要で、型が合っていればよい

	runner := di.ExportProvideSharedTxRunner(db)
	if runner == nil {
		t.Fatalf("expected non-nil tx runner from ExportProvideSharedTxRunner")
	}

	// stx.Runner インタフェースとして扱えることをコンパイル時に保証
	var _ stx.Runner = runner
}
