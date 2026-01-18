// Package request は、HTTP リクエストから受け取るデータ転送オブジェクト (DTO) を定義します。
// このパッケージでは、主に controller 層で Gin の Bind 機能によって入力を構造体にマッピングし、
// validation・正規化を行うための型を提供します。
package request

import (
	"strings"

	vo "resume/internal/usecase/valueobject"
)

// PatchProfileRequest は、プロフィール情報（氏名・生年月日・性別など）の
// 更新リクエストを表す構造体です。
//
// PATCH /fl/profile/personal のリクエストボディに対応し、
// JSON → DTO へのバインド時に Gin の binding タグによって
// 必須項目やオプショナル項目のバリデーションを行います。
//
// - `BirthDate` / `GenderID` / `Initial` は省略可能です。
// JSON にフィールドが存在しない場合は Set=false となり、ユースケース層で「更新しない」と判断されます。
// JSON に null が指定された場合は Set=true, Value=nil となり、「値をクリアする」として扱われます。
//
// Normalize メソッドによって、入力値の前後スペース除去などの軽微な正規化を行います。
type PatchProfileRequest struct {
	FamilyName     string                `json:"family_name" binding:"required"`
	GivenName      string                `json:"given_name" binding:"required"`
	FamilyNameKana string                `json:"family_name_kana" binding:"required"`
	GivenNameKana  string                `json:"given_name_kana" binding:"required"`
	BirthDate      vo.PatchValue[string] `json:"birth_date" binding:"omitempty"`
	GenderID       vo.PatchValue[uint8]  `json:"gender_id" binding:"omitempty"`
	Initial        vo.PatchValue[string] `json:"initial" binding:"omitempty"`
}

// Normalize はハンドラ層で Bind 後に呼ぶ
func (r *PatchProfileRequest) Normalize() {
	trim := func(p *string) {
		if p == nil {
			return
		}
		s := strings.TrimSpace(*p)
		if s == "" {
			*p = ""
		} else {
			*p = s
		}
	}

	trim(&r.FamilyName)
	trim(&r.GivenName)
	trim(&r.FamilyNameKana)
	trim(&r.GivenNameKana)

	//trim(r.BirthDate)
	//trim(r.GenderID)
	//trim(r.Initial)
}
