// Package field は、フィールドレベルで適用されるカスタムバリデーションルールを提供します。
// このファイルでは日本の郵便番号（jp_postal）に関する検証を定義しています。
package field

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
)

var reJPPostal = regexp.MustCompile(`^\d{3}-\d{4}$`) // 末尾は $

// jpPostal は、日本の郵便番号（NNN-NNNN 形式）を検証します。
func jpPostal(fl validator.FieldLevel) bool {
	val := strings.TrimSpace(fl.Field().String())
	if val == "" {
		return true
	}
	return reJPPostal.MatchString(val)
}

// registerJPPostal は、郵便番号用のカスタム検証関数 jp_postal を validator に登録します。
func registerJPPostal(v *validator.Validate) error {
	return v.RegisterValidation("jp_postal", jpPostal)
}

func init() { rules.RegisterField(registerJPPostal) }
