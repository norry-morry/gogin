// Package i18n は、国際化 (i18n) に関する値オブジェクト群を定義します。
// 言語バンドルやロケール情報などを管理します。
package i18n

// Locale は "ja" / "en" などの言語ロケールを表します。
// 値オブジェクトとして等価性と正規化を定義します。
type Locale struct {
	code string
}

// NewLocale はロケールコードを正規化して生成します。
func NewLocale(code string) Locale {
	// 必要に応じて大小文字変換やバリデーションを追加
	return Locale{code: code}
}

// Code は内部表現を返します。
func (l Locale) Code() string {
	return l.code
}

// Equal は値オブジェクトの等価比較です。
func (l Locale) Equal(other Locale) bool {
	return l.code == other.code
}

// String 実装（fmt.Printf等での出力用）
func (l Locale) String() string {
	return l.code
}
