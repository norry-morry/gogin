// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import "context"

// GetMessage は単一キーの翻訳を取得します。
// Translator 経由でキャッシュを参照し、該当キーが存在する場合は文字列を返します。
func (uc *Interactor) GetMessage(ctx context.Context, in GetMessageInput) (GetMessageOutput, error) {
	txt, ok := uc.tr.Translate(in.Locale, in.Key)
	return GetMessageOutput{Text: txt, OK: ok}, nil
}
