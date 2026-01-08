// Package jp は、日本ロケール固有のバリデーションルールを提供します。
// 例として、CountryCode が "JP" の場合のみ都道府県必須など、
// 日本向けの住所検証ロジックを含みます。
package jp

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// WrapWhenJP は CountryCode が "JP" の場合にのみ、
// 指定タグ（例: "jp_pref_required"）を適用する構造体レベルのバリデータを生成します。
func WrapWhenJP(tag string) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		cur := sl.Current()
		cc := cur.FieldByName("CountryCode")
		if !cc.IsValid() || cc.Kind() != reflect.String {
			return
		}
		if !strings.EqualFold(cc.String(), "JP") {
			return
		}
		pref := cur.FieldByName("AdministrativeArea")
		if !pref.IsValid() || pref.Kind() != reflect.String {
			return
		}
		if strings.TrimSpace(pref.String()) == "" {
			sl.ReportError(pref.Interface(), "administrativeArea", "AdministrativeArea", tag, "required")
		}
	}
}
