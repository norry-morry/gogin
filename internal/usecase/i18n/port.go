// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import vo "resume/internal/domain/valueobject/i18n"

// Translator は翻訳キャッシュなどを抽象化した読み取りポートです。
// 実装は adapter/gateway 側（CacheStore など）に配置します。
type Translator interface {
	// Translate は単一キーの翻訳を返します。
	Translate(loc vo.Locale, key vo.Key) (string, bool)

	// Bundle は prefix 前方一致で辞書エントリをまとめて返します。例: "ui.page.profile."
	Bundle(loc vo.Locale, prefix string) map[string]string
}
