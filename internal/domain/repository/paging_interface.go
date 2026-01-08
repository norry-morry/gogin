// Package repository は ページングのインフラ層を抽象化したもの
package repository

// BaseListSpec は、一覧取得（リスト）系ユースケースで使用される共通の検索条件インターフェースです。
// ページング、ソートなどの仕様を統一的に扱うために利用されます。
type BaseListSpec interface {
	// Offset は、検索結果の開始位置（スキップ件数）を返します。
	// 通常は (Page - 1) * Limit で算出されます。
	Offset() int

	// Limit は、取得件数の上限を返します。
	// -1 の場合は全件取得を意味します。
	Limit() int

	// SortCol は、ソート対象のカラム名を返します。
	// 例: "created_at", "id" など。
	SortCol() string

	// Order は、ソート順を返します。
	// "asc"（昇順）または "desc"（降順）のいずれかです。
	Order() string
}
