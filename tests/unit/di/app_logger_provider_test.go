// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルではアプリケーションロガー用の
// プロバイダのテストを行います。
package di_test

import (
	"testing"

	"resume/internal/config"
	"resume/internal/di"
)

// TestNewAppLogger_ReturnsLogger は ExportNewAppLogger が
// 非 nil の *slog.Logger を返すことを確認します。
func TestNewAppLogger_ReturnsLogger(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	// AppLogPath は空でもよい前提（実装側でハンドリングされる）
	logger := di.ExportNewAppLogger(cfg)
	if logger == nil {
		t.Fatalf("expected non-nil logger from ExportNewAppLogger")
	}
}
