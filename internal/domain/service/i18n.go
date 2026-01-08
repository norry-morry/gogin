// Package service は i18n に関連するドメインサービスを提供します。
package service

import vo "resume/internal/domain/valueobject/i18n"

// MergeBundles は、複数の辞書バンドルをマージして新しいバンドルを生成するドメインサービスです。
// 同一LocaleのBundleが複数ある場合、後の要素の値で上書きされます。
func MergeBundles(bundles ...*vo.Bundle) *vo.Bundle {
	if len(bundles) == 0 {
		return vo.NewBundle(vo.NewLocale(""), nil)
	}

	loc := bundles[0].Locale()
	out := map[string]string{}

	for _, b := range bundles {
		for k, v := range b.Messages() {
			out[k] = v // 後勝ち
		}
	}

	return vo.NewBundle(loc, out)
}
