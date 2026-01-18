// Package rulestruct は、構造体レベル/フィールドレベルで適用される
// 一覧系 DTO 向けのカスタムバリデーションルールを提供します。
// sort_key は、エンドポイントごとに許容されるソートキーを切り替える
// フィールド検証ルールです。
package rulestruct

import (
	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
)

// sortKeyValidator は、struct タグ `sort_key=identities` のように利用される
// カスタムフィールドバリデータです。
func sortKeyValidator(fl validator.FieldLevel) bool {
	param := fl.Param() // "identities" 等
	sortValue := fl.Field().String()

	// 未指定は OK（Normalize でデフォルトが入る設計）
	if sortValue == "" {
		return true
	}

	allowed, ok := rules.SortKeyAllowed[param]
	if !ok {
		// 未定義の sort_key=xxx はエラー（運用・開発者向け）
		return false
	}

	for _, v := range allowed {
		if v == sortValue {
			return true
		}
	}
	return false
}

// registerSortKey は validator に sort_key バリデータを登録します。
func registerSortKey(v *validator.Validate) error {
	return v.RegisterValidation("sort_key", sortKeyValidator)
}

// init で rules に登録 → 最終的に ApplyAll() により実体登録される
func init() { rules.RegisterField(registerSortKey) }
