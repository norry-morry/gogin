// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
// Package db は GORM 用の SQL ログ出力を含むロギングユーティリティを提供します。
package db

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger は slog を用いた GORM の logger.Interface 実装です。
type GormLogger struct {
	base           *slog.Logger
	level          logger.LogLevel
	slow           time.Duration
	ignoreNotFound bool
}

// NewGormLogger は GormLogger を生成します。
func NewGormLogger(base *slog.Logger, level logger.LogLevel, slow time.Duration, ignoreNotFound bool) *GormLogger {
	return &GormLogger{base: base, level: level, slow: slow, ignoreNotFound: ignoreNotFound}
}

// LogMode はログレベルを切り替えた新しいロガーを返します。
func (g *GormLogger) LogMode(l logger.LogLevel) logger.Interface { cp := *g; cp.level = l; return &cp }

// Info は GORM の情報ログを出力します。
func (g *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if g.level >= logger.Info {
		g.base.InfoContext(ctx, msg, data...)
	}
}

// Warn は GORM の警告ログを出力します。
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if g.level >= logger.Warn {
		g.base.WarnContext(ctx, msg, data...)
	}
}

// Error は GORM のエラーログを出力します。
func (g *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if g.level >= logger.Error {
		g.base.ErrorContext(ctx, msg, data...)
	}
}

// Trace は SQL の実行結果および所要時間等を出力します。
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.level == logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil:
		if g.ignoreNotFound && errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		g.base.ErrorContext(ctx, "sql.error", "err", err, "elapsed", elapsed, "rows", rows, "sql", sql)
	case g.slow > 0 && elapsed > g.slow:
		g.base.WarnContext(ctx, "sql.slow", "elapsed", elapsed, "rows", rows, "sql", sql)
	default:
		if g.level >= logger.Info {
			g.base.InfoContext(ctx, "sql.trace", "elapsed", elapsed, "rows", rows, "sql", sql)
		}
	}
}
