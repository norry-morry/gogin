// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import "context"

// ListBundle は prefix 前方一致で翻訳辞書をまとめて取得します。
// Translator 経由でキャッシュから対応するキー群を抽出します。
func (uc *Interactor) ListBundle(ctx context.Context, in ListBundleInput) (ListBundleOutput, error) {
	data := uc.tr.Bundle(in.Locale, in.Prefix)
	return ListBundleOutput{Data: data}, nil
}
