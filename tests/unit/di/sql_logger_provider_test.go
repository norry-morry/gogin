// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルでは SQL ログ／GORM ロガー用の
// プロバイダのテストを行います。
package di_test

import (
	"testing"

	"gorm.io/gorm/logger"

	"resume/internal/config"
	"resume/internal/di"
)

// TestNewSQLLogger_ReturnsSQLLog は ExportNewSQLLogger が
// 非 nil の SQLLog と、その中のロガーを返すことを確認します。
func TestNewSQLLogger_ReturnsSQLLog(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	sqlLog := di.ExportNewSQLLogger(cfg)
	if sqlLog == nil {
		t.Fatalf("expected non-nil SQLLog from ExportNewSQLLogger")
	}
	if sqlLog.L == nil {
		t.Fatalf("expected SQLLog.L to be non-nil")
	}
}

// TestProvideGormLogger_ReturnsGormLogger は ExportProvideGormLogger が
// gorm の logger.Interface を実装するロガーを返すことを確認します。
func TestProvideGormLogger_ReturnsGormLogger(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	sqlLog := di.ExportNewSQLLogger(cfg)

	gl := di.ExportProvideGormLogger(sqlLog)
	if gl == nil {
		t.Fatalf("expected non-nil GormLogger from ExportProvideGormLogger")
	}

	// gorm の logger.Interface として扱えることをコンパイル時に保証
	var _ logger.Interface = gl
}
