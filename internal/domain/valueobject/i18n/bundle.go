// Package i18n は、国際化 (i18n) に関する値オブジェクト群を定義します。
// 言語バンドルやロケール情報などを管理します。
package i18n

// Bundle は 1 ロケール分の翻訳辞書を表す値オブジェクトです。
// 原則イミュータブルとして扱います。
type Bundle struct {
	locale   Locale
	messages map[string]string // Key.Value() を key にして格納
}

// NewBundle は新しい Bundle を構築します。
func NewBundle(locale Locale, messages map[string]string) *Bundle {
	cp := make(map[string]string, len(messages))
	for k, v := range messages {
		cp[k] = v
	}
	return &Bundle{
		locale:   locale,
		messages: cp,
	}
}

// Locale はこの Bundle のロケールを返します。
func (b *Bundle) Locale() Locale {
	return b.locale
}

// Get はキーに対応するメッセージを取得します。
func (b *Bundle) Get(k Key) (string, bool) {
	v, ok := b.messages[k.Value()]
	return v, ok
}

// Messages は全メッセージをコピーで返します（安全のため）。
func (b *Bundle) Messages() map[string]string {
	cp := make(map[string]string, len(b.messages))
	for k, v := range b.messages {
		cp[k] = v
	}
	return cp
}
