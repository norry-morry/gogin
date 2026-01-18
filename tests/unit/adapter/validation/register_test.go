// Package validation_test は internal/adapter/validation パッケージに定義された
// validator 初期化処理およびカスタムルール登録処理に対するユニットテストを提供します。
package validation_test

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation"
)

// fakeValidator は gin/binding.StructValidator を満たすテスト用実装です。
// Engine() が *validator.Validate 以外を返すパターンを再現するために使用します。
type fakeValidator struct{}

func (f *fakeValidator) ValidateStruct(obj interface{}) error {
	return nil
}

func (f *fakeValidator) Engine() interface{} {
	return struct{}{}
}

// TestMultiRegisterAll_PanicsWhenEngineIsNotValidator は、
// binding.Validator.Engine() の戻り値が *validator.Validate でない場合に
// MultiRegisterAll が panic することを確認します。
func TestMultiRegisterAll_PanicsWhenEngineIsNotValidator(t *testing.T) {
	// 並列にすると binding.Validator の書き換えが競合するので t.Parallel() は使わない

	origValidator := binding.Validator
	defer func() {
		binding.Validator = origValidator
		if r := recover(); r == nil {
			t.Errorf("expected panic when Engine is not *validator.Validate, but no panic occurred")
		}
	}()

	// Engine() が *validator.Validate ではない StructValidator を差し込む
	binding.Validator = &fakeValidator{}

	validation.MultiRegisterAll()
}

// TestMultiRegisterAll_SucceedsWithDefaultValidator は、
// validator.Validate を内部に持つ StructValidator を使った場合に
// MultiRegisterAll が panic せず完走することを確認します。
func TestMultiRegisterAll_SucceedsWithDefaultValidator(t *testing.T) {
	// ここも binding.Validator を書き換えるため t.Parallel() は使わない

	origValidator := binding.Validator
	defer func() {
		binding.Validator = origValidator
		if r := recover(); r != nil {
			t.Fatalf("MultiRegisterAll() panicked unexpectedly: %v", r)
		}
	}()

	// テスト用の「安全な」 StructValidator を明示的にセットする
	binding.Validator = NewDefaultStructValidatorForTest()

	validation.MultiRegisterAll()
}

// NewDefaultStructValidatorForTest は *validator.Validate を内部に持つ
// gin/binding.StructValidator の簡易実装をテスト用に提供します。
func NewDefaultStructValidatorForTest() binding.StructValidator {
	return &defaultTestValidator{validate: validator.New()}
}

type defaultTestValidator struct {
	validate *validator.Validate
}

func (v *defaultTestValidator) ValidateStruct(obj interface{}) error {
	return v.validate.Struct(obj)
}

func (v *defaultTestValidator) Engine() interface{} {
	return v.validate
}
