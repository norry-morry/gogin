// Package presenter は、ユースケース層の出力データを API レスポンス用 DTO に変換する責務を持ちます。
// ドメインロジックや永続化層の構造をクライアントへ直接公開しないための中間変換レイヤです。
package presenter

import (
	"resume/internal/adapter/http/dto/response"
	"resume/internal/shared/util"
	"resume/internal/usecase/profile"
)

// PresentUserProfile は、ユースケース層の UserProfileOutput を HTTP レスポンス用の
// UserProfileResponse DTO に変換します。
// 生年月日のフォーマット変換や年齢グループ算出など、表示専用の補正処理を行います。
func PresentUserProfile(out profile.UserProfileOutput) response.UserProfileResponse {
	var birthStr *string
	if out.BirthDate != nil {
		birthStr = util.FormatDatePtr(out.BirthDate)
	}

	var ageGroup *int

	if out.Age != nil {
		g := util.CalcAgeGroup(*out.Age) // 25 -> 20, 38 -> 30 など
		if g > 0 {
			ageGroup = &g
		}
	}

	return response.UserProfileResponse{
		UserID:         out.UserID,
		FamilyName:     out.FamilyName,
		GivenName:      out.GivenName,
		FamilyNameKana: out.FamilyNameKana,
		GivenNameKana:  out.GivenNameKana,
		LegalName:      out.LegalName,
		LegalNameKana:  out.LegalNameKana,
		BirthDate:      birthStr,
		Age:            out.Age,
		AgeGroup:       ageGroup,
		GenderID:       out.GenderID,
		GenderLabelKey: out.GenderLabelKey,
		Initial:        out.Initial,
	}
}
