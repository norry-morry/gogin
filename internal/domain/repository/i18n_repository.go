// Package repository は 多言語化辞書のインフラ層抽象
package repository

import vo "resume/internal/domain/valueobject/i18n"

// DictionaryRepository は、i18n辞書データをロードするリポジトリポートです。
// 実装は adapter/gateway 側（FS/db/S3など）に配置します。
type DictionaryRepository interface {
	// LoadAll 全ロケールの辞書をロードします。
	// 例: map["ja"] => *Bundle, map["en"] => *Bundle
	LoadAll() (map[vo.Locale]*vo.Bundle, error)
}
