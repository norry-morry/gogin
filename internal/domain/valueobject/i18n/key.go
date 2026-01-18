// Package i18n は、国際化 (i18n) に関する値オブジェクト群を定義します。
// 言語バンドルやロケール情報などを管理します。
package i18n

// Key は "ui.page.profile.title" のような翻訳キーを表す値オブジェクトです。
type Key struct {
	value string
}

// NewKey はキー文字列を生成します。
func NewKey(v string) Key {
	// ここで形式チェック（例: ドット区切り強制など）を加えても良い
	return Key{value: v}
}

// Value は内部値を返します。
func (k Key) Value() string {
	return k.value
}

// Equal は等価比較を行います。
func (k Key) Equal(other Key) bool {
	return k.value == other.value
}

// String はデバッグ等での表示用
func (k Key) String() string {
	return k.value
}
