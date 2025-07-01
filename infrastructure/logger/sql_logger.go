package logger

import (
	"context"
	"fmt"
	gormlogger "gorm.io/gorm/logger"
	"log/slog"
	"time"
)

type SlogGormLogger struct {
	logger        *slog.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func NewSlogGormLogger(logger *slog.Logger, level gormlogger.LogLevel) gormlogger.Interface {
	return &SlogGormLogger{
		logger:        logger,
		LogLevel:      level,
		SlowThreshold: 200 * time.Millisecond, // 遅いクエリの閾値
	}
}

func (l *SlogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []slog.Attr{
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		fields = append(fields, slog.Any("error", err))
		//l.logger.Error("gorm error", fields) // ← ここ
		args := make([]any, len(fields))
		for i, f := range fields {
			args[i] = f
		}
		l.logger.Error("gorm error", slog.Group("fields", args...))
	case elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		fields = append(fields, slog.String("warning", "slow query"))
		//l.logger.Warn("gorm slow query", fields) // ← ここ
		args := make([]any, len(fields))
		for i, f := range fields {
			args[i] = f
		}
		l.logger.Warn("gorm slow query", slog.Group("fields", args...))
	case l.LogLevel >= gormlogger.Info:
		//l.logger.Info("gorm query", fields) // ← ここ
		args := make([]any, len(fields))
		for i, f := range fields {
			args[i] = f
		}
		l.logger.Warn("gorm query", slog.Group("fields", args...))
	}
}
