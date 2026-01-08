// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルでは認証用 Clock プロバイダの
// テストを行います。
package di_test

import (
	"testing"

	"resume/internal/di"
)

// TestProvideAuthClock_CurrentlyNil は ExportProvideAuthClock が
// 現状 nil を返す実装であることを確認します。
// 将来、具体的な Clock 実装へ差し替える場合はこのテストを更新します。
func TestProvideAuthClock_CurrentlyNil(t *testing.T) {
	t.Parallel()

	clk := di.ExportProvideAuthClock()
	if clk != nil {
		t.Errorf("expected nil clock from ExportProvideAuthClock (current design), got=%v", clk)
	}
}
