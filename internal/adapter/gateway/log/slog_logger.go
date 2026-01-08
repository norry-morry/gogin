// Package log Package config provides logger configuration settings.
package log

import (
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	shared "resume/internal/shared/log"
)

// InitRotatingLogger は指定ファイルにローテーション付きで出力する slog ロガーを初期化します
func InitRotatingLogger(filePath string) *slog.Logger {
	stdout := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	)
	file := slog.NewJSONHandler(
		&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    50,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	)
	// shared の MultiHandler があるならそれを利用。無ければ自前で MultiWriter 的に。
	mh := shared.NewMultiHandler(stdout, file)
	l := slog.New(mh)
	slog.SetDefault(l)
	return l
}
