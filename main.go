package main

import (
	"log"
	"log/slog"
	"resume/config"
	"resume/di"
)

func main() {
	// 環境変数の読込
	cfg := config.Load()

	// リクエストログを含む各種ロガーの設定
	logger := config.SetupLogger(cfg.LogPath)
	slog.SetDefault(logger)

	// DIコンテナでrouterを構成
	r, err := di.InitApp(cfg)
	if err != nil {
		//log.Fatalf("di.InitApp err: %v", err)
		logger.Error("failed to initialize app", slog.Any("error", err))
		log.Fatal(err)
	}

	logger.Info("Hello golang from docker with air!")
	//r.Run(cfg.AppIP + ":" + cfg.AppPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := r.Run(cfg.AppIP + ":" + cfg.AppPort); err != nil {
		logger.Error("server error", slog.Any("error", err))
	}
}
