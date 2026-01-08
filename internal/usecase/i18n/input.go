// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import vo "resume/internal/domain/valueobject/i18n"

// GetMessageInput は単一キーの翻訳を取得するユースケースへの入力DTOです。
type GetMessageInput struct {
	Locale vo.Locale
	Key    vo.Key
}

// ListBundleInput は prefix で翻訳辞書をまとめて取得するユースケースへの入力DTOです。
type ListBundleInput struct {
	Locale vo.Locale
	Prefix string
}
