// Package bootstrap はプロセス起動時の環境初期化を扱います。
package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadDotEnv は開発環境で .env を読み込みます。
func LoadDotEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" || env == "local" || env == "dev" || env == "develop" || env == "development" || env == "stg" || env == "staging" {
		// CIや本番では env_file / 環境変数経由で注入される想定。失敗は致命でない。
		//nolint:errcheck
		_ = godotenv.Load()
	}
}
