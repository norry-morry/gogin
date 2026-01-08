// Package util は pointerに関する軽ユーティリティです
package util

// Clone は与えられたポインタの値をコピーして新しいポインタを返します。
// nil の場合は nil を返します。
func Clone[T any](in *T) *T {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

// ToPtr は値をポインタ化します。
// 例: util.ToPtr("foo") → *string("foo")
func ToPtr[T any](v T) *T {
	return &v
}

// Deref はポインタを値に戻します。nil の場合はゼロ値を返します。
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
