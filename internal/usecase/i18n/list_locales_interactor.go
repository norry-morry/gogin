// Package i18n は、多言語辞書のロードおよび翻訳ロジックを管理します。
package i18n

import vo "resume/internal/domain/valueobject/i18n"

// ListLocales はキャッシュまたはリポジトリから利用可能ロケール一覧を取得します。
func (uc *Interactor) ListLocales() []vo.Locale {
	// まずキャッシュに存在する場合はそれを返す
	if locales := uc.cache.Locales(); len(locales) > 0 {
		return locales
	}

	// fallback: repoからロードしてキー一覧をLocaleに変換
	all, err := uc.repo.LoadAll()
	if err != nil {
		return []vo.Locale{}
	}
	out := make([]vo.Locale, 0, len(all))
	for loc := range all {
		out = append(out, loc)
	}
	return out
}
