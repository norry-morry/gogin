package field

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
)

func noEdgeSpaces(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return s == strings.TrimSpace(s)
}

// registerNoEdgeSpaces は、前後空白検知用のカスタム検証関数 no_edge_spaces を validator に登録します。
func registerNoEdgeSpaces(v *validator.Validate) error {
	return v.RegisterValidation("no_edge_spaces", noEdgeSpaces)
}

func init() {
	rules.RegisterField(registerNoEdgeSpaces)
}
