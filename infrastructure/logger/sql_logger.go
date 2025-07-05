// Package logger は GORM 用の SQL ログ出力を含むロギングユーティリティを提供します。
package logger

import (
	"context"
	"time"

	"log/slog"

	"gorm.io/gorm/logger"
)

// SlogGormLogger は slog を使用した GORM のロガー実装です。
type SlogGormLogger struct {
	logger *slog.Logger
	level  logger.LogLevel
}

// NewSlogGormLogger は GORM のロガーインターフェースを実装する新しい SlogGormLogger を作成して返します。
func NewSlogGormLogger(logger *slog.Logger, level logger.LogLevel) logger.Interface {
	return &SlogGormLogger{
		logger: logger,
		level:  level,
	}
}

// LogMode は SlogGormLogger のログレベルを設定し、更新されたロガーを返します。
func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

// Info は指定されたコンテキスト、メッセージ、追加データを使用して情報ログを出力します。
func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.InfoContext(ctx, msg, data...)
}

// Warn は指定されたコンテキスト、メッセージ、追加データを使用して警告ログを出力します。
func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WarnContext(ctx, msg, data...)
}

// Error は指定されたコンテキスト、メッセージ、追加データを使用してエラーログを出力します。
func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ErrorContext(ctx, msg, data...)
}

// Trace はデータベース操作の詳細（実行時間、実行SQL、影響行数、エラー）をログ出力します。
func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.logger.ErrorContext(ctx, "SQL エラー", "error", err, "duration", elapsed, "rows", rows, "sql", sql)
	} else {
		l.logger.InfoContext(ctx, "SQL トレース", "duration", elapsed, "rows", rows, "sql", sql)
	}
}
