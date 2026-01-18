package rules

import "github.com/go-playground/validator/v10"

type fieldRegistrar func(*validator.Validate) error
type structRegistrar func(*validator.Validate)

var fieldRegs []fieldRegistrar
var structRegs []structRegistrar

// RegisterField は、フィールドレベルのカスタム検証関数を登録キューに追加します。
// 実際の登録は ApplyAll() 実行時に行われます。
func RegisterField(f fieldRegistrar) { fieldRegs = append(fieldRegs, f) }

// RegisterStruct は、構造体レベルのカスタム検証関数を登録キューに追加します。
// 実際の登録は ApplyAll() 実行時に行われます。
func RegisterStruct(f structRegistrar) { structRegs = append(structRegs, f) }

// ApplyAll は、これまで RegisterField / RegisterStruct により登録された
// すべてのカスタム検証関数を validator に適用します。
func ApplyAll(v *validator.Validate) error {
	for _, f := range fieldRegs {
		if err := f(v); err != nil {
			return err
		}
	}
	for _, f := range structRegs {
		f(v)
	}
	return nil
}
