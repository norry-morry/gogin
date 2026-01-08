// Package main は resume アプリケーションのエントリーポイントです。
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"resume/internal/bootstrap"
	"resume/internal/di"
	"resume/internal/shared/util"
)

func main() {
	// 1, .envの読込
	bootstrap.LoadDotEnv()

	// 2. validatorの登録
	bootstrap.SetupValidation()

	// 3. 機微情報マスキング設定
	setupSensitiveMasker()

	// 4. i18n 初期化 + DI 経由で router 組み立て
	engine, err := di.InitAPI("internal/shared/i18n/locales")
	if err != nil {
		log.Fatal(err)
	}

	// 5. アプリケーション起動
	//if err := engine.Run(); err != nil { // 既存の起動方法に合わせて
	//	log.Fatal(err)
	//}
	// 5. HTTP サーバ起動（Graceful Shutdown 付き）
	if err := runHTTPServer(engine); err != nil {
		log.Fatal(err)
	}
}

// 例: 起動時に一度だけ設定ファイルから読み込み、グローバルに差し替える
func setupSensitiveMasker() {
	path := os.Getenv("SENSITIVE_CONFIG")
	if path == "" {
		// 実行ファイル基準（相対パスのズレを避ける）
		exe, err := os.Executable()
		if err != nil {
			slog.Warn("sensitive: cannot get executable path", "err", err)
			return
		}
		base := filepath.Dir(exe)
		path = filepath.Join(base, "configs", "logging.yaml")
	}

	cfg, err := util.LoadSensitiveConfig(path)
	if err != nil {
		slog.Warn("sensitive: load failed, use noop", "path", path, "err", err)
		return
	}
	m, err := util.NewMasker(cfg)
	if err != nil {
		slog.Warn("sensitive: build failed, use noop", "err", err)
		return
	}
	util.SetGlobalMasker(m)

	// ✅ 起動時セルフテスト（ログで “本当に効いてる” を確認）
	test := map[string]any{
		"address_line1": "0123456789012",
		"email":         "foo@example.com",
		"password":      "secret",
	}
	masked := m.MaskMapShallow(test)
	slog.Info("sensitive: enabled",
		"config_path", path,
		"sample_masked", masked,
	)
}

// app.Run() 相当を main パッケージ側に移植
func runHTTPServer(handler http.Handler) error {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// 起動
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏳ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("✅ Server gracefully stopped")
	return nil
}
