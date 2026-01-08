package field

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
	"resume/internal/shared/jpstring"
)

func jpPref(fl validator.FieldLevel) bool {
	val := strings.TrimSpace(fl.Field().String())
	if val == "" {
		return true
	} // omitempty想定
	return jpstring.IsValidPrefecture(val)
}

// registerJPPref は、都道府県用のカスタム検証関数 jp_pref を validator に登録します。
func registerJPPref(v *validator.Validate) error {
	return v.RegisterValidation("jp_pref", jpPref)
}

func init() { rules.RegisterField(registerJPPref) }
