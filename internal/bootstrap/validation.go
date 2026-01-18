// Package bootstrap はプロセス起動時のバリデーションルール初期化を扱います。
package bootstrap

import (
	"log"

	"resume/internal/adapter/validation"
)

// SetupValidation はアプリケーション起動時に validator の初期設定を行います。
// カスタム検証ルールや構造体レベルのルールを一括で登録します。
func SetupValidation() {
	log.Println("[bootstrap] registering validator rules...")
	validation.MultiRegisterAll()
}
