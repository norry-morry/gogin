package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"log/slog"
	"os"
	"resume/Utility"
)

func SetupLogger() *slog.Logger {
	log.Println(os.Getenv("LOGGER_PATH"))
	logFile := &lumberjack.Logger{
		Filename:   os.Getenv("LOGGER_PATH"),
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
