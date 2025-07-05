// Package config provides logger configuration settings.
package config

import (
	"log/slog"
	"os"
	"resume/utility"

	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger は指定されたファイルパスにログを出力する slog.Logger を初期化して返します。
func SetupLogger(filePath string) *slog.Logger {

	logFile := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	// 標準出力
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	// ファイル出力
	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	multiHandler := utility.NewMultiHandler(stdoutHandler, fileHandler)

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	return logger
}
