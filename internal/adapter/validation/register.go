// Package validation は、validator の生成とカスタムルール登録を提供します。
package validation

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	req "resume/internal/adapter/http/dto/request"
	"resume/internal/adapter/validation/rules"
	_ "resume/internal/adapter/validation/rules/field" // revive:disable:blank-imports - init() でフィールド単位のカスタム検証を登録するための副作用インポート
	jplocale "resume/internal/adapter/validation/rules/locale/jp"
	rulestruct "resume/internal/adapter/validation/rules/struct"
)

// MultiRegisterAll は、Gin の binding.Validator に紐づいた *validator.Validate に対し、
// すべてのカスタム検証（フィールド/構造体/ロケール依存）をまとめて登録します。
// 既に binding.Validator が *validator.Validate でない場合は panic します。
func MultiRegisterAll() {
	eng, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Panic("gin validator engin not *validator.Validate")
	}

	if err := rules.ApplyAll(eng); err != nil {
		log.Panic(err)
	}

	// 構造体(汎用バリデータ)
	eng.RegisterStructValidation(rulestruct.CoordPair, req.CreateAddressRequest{})

	// 構造体(条件: JPの時だけ)
	eng.RegisterStructValidation(jplocale.WrapWhenJP("jp_pref_required"), req.CreateAddressRequest{})
}
