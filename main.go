// Package main は resume アプリケーションのエントリーポイントです。
package main

import (
	"log"
	"log/slog"
	"resume/config"
	"resume/di"
)

func main() {
	slog.Debug("Hello")
	// 環境変数の読込
	cfg := config.Load()

	// リクエストログを含む各種ロガーの設定
	appLogger := config.SetupLogger(cfg.AppLogPath)
	slog.SetDefault(appLogger)

	sqlLogger := config.SetupLogger(cfg.SQLLogPath)

	// DIコンテナでrouterを構成
	r, err := di.InitApp(cfg, sqlLogger)
	if err != nil {
		//log.Fatalf("di.InitApp err: %v", err)
		appLogger.Error("failed to initialize app", slog.Any("error", err))
		log.Fatal(err)
	}

	appLogger.Info("Hello golang from docker with air!")
	//r.Run(cfg.AppIP + ":" + cfg.AppPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := r.Run(cfg.AppIP + ":" + cfg.AppPort); err != nil {
		appLogger.Error("server error", slog.Any("error", err))
	}
}
