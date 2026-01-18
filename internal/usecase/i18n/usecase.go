// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import (
	"context"

	vo "resume/internal/domain/valueobject/i18n"
)

// Usecase は i18n のユースケース群を表すインターフェイスです。
// 翻訳取得・バンドル一覧取得・辞書リロードの操作を統合しています。
type Usecase interface {
	// GetMessage は単一キーの翻訳を返します。
	GetMessage(ctx context.Context, in GetMessageInput) (GetMessageOutput, error)

	// ListBundle は prefix で辞書エントリをまとめて取得します。
	ListBundle(ctx context.Context, in ListBundleInput) (ListBundleOutput, error)

	// Reload は外部リポジトリから辞書を再読込し、キャッシュへ反映します。
	Reload(ctx context.Context) error

	ListLocales() []vo.Locale
}
