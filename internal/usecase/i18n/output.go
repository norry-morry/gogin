// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

// GetMessageOutput は単一キーの翻訳結果を表します。
type GetMessageOutput struct {
	Text string
	OK   bool
}

// ListBundleOutput は prefix で抽出された翻訳辞書の結果を表します。
type ListBundleOutput struct {
	Data map[string]string
}
