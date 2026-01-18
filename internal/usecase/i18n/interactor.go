// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

// Interactor は i18n ユースケースの実行体です。
// Translator, DictionaryRepository, CacheStore を協調させて翻訳やリロードを行います。
type Interactor struct {
	tr    Translator
	repo  DicRepo
	cache CacheStore
}

// New は i18nユースケースの実装（Interactor）を生成します。
func New(
	tr Translator,
	repo DicRepo,
	cache CacheStore,
) Usecase {
	return &Interactor{
		tr:    tr,
		repo:  repo,
		cache: cache,
	}
}
