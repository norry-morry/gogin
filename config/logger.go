package config

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

func SetupLogger(filePath string) *slog.Logger {
	logFile := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	// 標準出力とファイル出力の両方
	writer := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelInfo, // or slog.LevelDebug
		AddSource: true,
	})
	return slog.New(handler)
}
