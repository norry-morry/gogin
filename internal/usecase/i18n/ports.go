// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import (
	repo "resume/internal/domain/repository"
	vo "resume/internal/domain/valueobject/i18n"
)

// DicRepo は i18n用の DictionaryRepository 型エイリアスです。
// 永続層から辞書データを取得するためのポートを表します。
type DicRepo = repo.DictionaryRepository

// CacheStore はキャッシュ層を抽象化したポートです。
// SwapAll により全ロケール分の辞書を一括で差し替えます。
type CacheStore interface {
	SwapAll(all map[vo.Locale]*vo.Bundle) error
	Get(loc vo.Locale, key vo.Key) (string, bool)
	Set(loc vo.Locale, key vo.Key, value string)
	Clear()
	Locales() []vo.Locale
}
