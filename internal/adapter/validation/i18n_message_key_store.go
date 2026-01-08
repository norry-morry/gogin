// Package validation は、ドメインロジックで発生した検証エラーを
// ユーザー向けの多言語メッセージへ変換するためのアダプタ群を提供します。
package validation

import (
	"resume/internal/adapter/validation/rules"
	vo "resume/internal/domain/valueobject/i18n"
	uci18n "resume/internal/usecase/i18n"
)

//
// I18nMessageKeyStore
// --------------------
// validation.rules / validation.custom / validation.fieldspecific
// などの翻訳キーが「辞書内に存在するか」を確認するための
// バックエンド側のキーリゾルバ（存在チェック専用）。
//
// MapValidationErrors（rules パッケージ内）から呼び出され、
// その中で「fieldspecific → custom → rules」の優先度で
// 最終的な Code（翻訳キー）を決定するために使用されます。
//
// ※この構造体は Translator のインスタンスを内部に保持し、
//   実際の YAML 辞書（validation.yaml 等）にアクセスします。
//   値自体（翻訳文字列）は使わず、キーが存在するかの bool だけを返します。
//

// I18nMessageKeyStore は、i18n Translator を利用して
// 「指定キーが存在するかどうか」を判定する MessageKeyStore 実装です。
type I18nMessageKeyStore struct {
	tr     uci18n.Translator // Translator（CacheStoreImpl が実装している）
	locale vo.Locale         // 対象ロケール（例: ja, en）
}

// NewI18nMessageKeyStore は、Translator とロケールを受け取って
// rules.MessageKeyStore を構築します。
// 通常は DI (wire.go) で Translator 初期化後に呼び出され、
// rules.SetMessageKeyStore(...) へ登録されます。
func NewI18nMessageKeyStore(tr uci18n.Translator, locale vo.Locale) rules.MessageKeyStore {
	return &I18nMessageKeyStore{
		tr:     tr,
		locale: locale,
	}
}

// HasKey は、指定キーが i18n 辞書内に存在するかを判定します。
// ここでは翻訳結果の文字列値は使用せず、存在チェックの bool のみ返します。
func (s *I18nMessageKeyStore) HasKey(key string) bool {
	_, ok := s.tr.Translate(s.locale, vo.NewKey(key))
	return ok
}

//
// 使い方イメージ（wire.go または main.go 側）
// ---------------------------------------------------------
//  uc, err := provideI18nUsecase(...)
//  if err != nil { ... }
//
//  // Translator は CacheStoreImpl として組まれている想定
//  keyStore := validation.NewI18nMessageKeyStore(tr, vo.NewLocale("ja"))
//  rules.SetMessageKeyStore(keyStore)
// ---------------------------------------------------------
//
// これにより、MapValidationErrors 内で hasValidationKey() を呼ぶと
// バックエンドの i18n キャッシュから「キー存在チェック」が走り、
// fieldspecific > custom > rules の優先度が自動的に適用されます。
//
