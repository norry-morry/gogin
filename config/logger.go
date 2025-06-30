package config

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"resume/Utility"
)

func SetupLogger(p string) *slog.Logger {
	//log.Println(p)
	//log.Println(os.Getenv("LOGGER_PATH"))
	logFile := &lumberjack.Logger{
		Filename:   p,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	multiHandler := Utility.NewMultiHandler(stdoutHandler, fileHandler)
	return slog.New(multiHandler)
}
