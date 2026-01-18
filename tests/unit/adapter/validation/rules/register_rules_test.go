// Package rules_test は internal/adapter/validation/rules パッケージに定義された
// カスタムバリデーションルール登録用ユーティリティのユニットテストを提供します。
package rules_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/validation/rules"
)

// TestApplyAll_RegistersFieldAndStruct は RegisterField / RegisterStruct に登録した
// 関数が ApplyAll によって validator に適用されることを確認します。
func TestApplyAll_RegistersFieldAndStruct(t *testing.T) {
	// グローバルな登録キューを使うため、並列実行は避ける（t.Parallel は使わない）
	v := validator.New()

	fieldCalled := false
	structCalled := false
	fieldValidatorMismatched := false
	structValidatorMismatched := false

	// フィールドレベル登録関数
	rules.RegisterField(func(got *validator.Validate) error {
		fieldCalled = true
		if got != v {
			// ここでは testing.T を触らずフラグだけ立てる
			fieldValidatorMismatched = true
		}
		return nil
	})

	// 構造体レベル登録関数
	rules.RegisterStruct(func(got *validator.Validate) {
		structCalled = true
		if got != v {
			structValidatorMismatched = true
		}
	})

	if err := rules.ApplyAll(v); err != nil {
		t.Fatalf("ApplyAll returned error: %v", err)
	}

	if !fieldCalled {
		t.Errorf("expected field registrar to be called")
	}
	if !structCalled {
		t.Errorf("expected struct registrar to be called")
	}
	if fieldValidatorMismatched {
		t.Errorf("field registrar received unexpected validator instance")
	}
	if structValidatorMismatched {
		t.Errorf("struct registrar received unexpected validator instance")
	}
}

// TestApplyAll_ReturnsErrorWhenFieldRegistrarFails は、
// RegisterField に登録した関数がエラーを返した場合に
// ApplyAll がそのエラーを呼び出し元へ返すことを確認します。
func TestApplyAll_ReturnsErrorWhenFieldRegistrarFails(t *testing.T) {
	// ここでもグローバルキューを使うので t.Parallel は使わない
	v := validator.New()

	expectedErr := errors.New("field registrar failed")

	// エラーを返すフィールドレベル登録関数
	rules.RegisterField(func(_ *validator.Validate) error {
		return expectedErr
	})

	err := rules.ApplyAll(v)
	if err == nil {
		t.Fatalf("expected error from ApplyAll, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
