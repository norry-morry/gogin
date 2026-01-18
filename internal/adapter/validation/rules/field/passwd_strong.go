package field

import (
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
)

func passwdStrong(fl validator.FieldLevel) bool {
	min := 12
	if p := fl.Param(); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			min = n
		}
	}
	s := fl.Field().String()
	if len([]rune(s)) < min {
		return false
	}

	var lower, upper, digit, symbol bool
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			lower = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsDigit(r):
			digit = true
		default:
			symbol = true
		}
	}
	cnt := 0
	for _, b := range []bool{lower, upper, digit, symbol} {
		if b {
			cnt++
		}
	}
	return cnt >= 3
}

// registerPasswdStrong は、パスワード強度検知用のカスタム検証関数 passwd_strong を validator に登録します。
func registerPasswdStrong(v *validator.Validate) error {
	return v.RegisterValidation("passwd_strong", passwdStrong)
}

func init() { rules.RegisterField(registerPasswdStrong) }
