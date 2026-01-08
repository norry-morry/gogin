// Package i18n は i18n（多言語化）関連のユースケースを提供します。
package i18n

import "context"

// Reload は永続層（DictionaryRepository）から全ロケール辞書を再取得し、
// CacheStore へ SwapAll して翻訳キャッシュを更新します。
func (uc *Interactor) Reload(ctx context.Context) error {
	all, err := uc.repo.LoadAll()
	if err != nil {
		return err
	}
	return uc.cache.SwapAll(all)
}
