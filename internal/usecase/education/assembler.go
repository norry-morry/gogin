// Package education は ユーザーの学歴情報に関するユースケースを提供します。
// このファイルは、ドメインエンティティを出力用データ構造やレスポンスDTOに変換するアセンブラ定義を含みます。
package education

import (
	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
	"resume/internal/shared/util"
)

func educationStatusKey(educationStatusID uint64) string {
	var key string
	switch educationStatusID {
	case entity.EducationStatusEnrolled:
		key = entity.EducationStatusCodeEnrolled
	case entity.EducationStatusLeaveOfAbsence:
		key = entity.EducationStatusCodeLeaveOfAbsence
	case entity.EducationStatusGraduated:
		key = entity.EducationStatusCodeGraduated
	case entity.EducationStatusCompleted:
		key = entity.EducationStatusCodeCompleted
	case entity.EducationStatusGraduationProspect:
		key = entity.EducationStatusCodeGraduationProspect
	case entity.EducationStatusWithdrawn:
		key = entity.EducationStatusCodeWithdrawn
	case entity.EducationStatusExpelled:
		key = entity.EducationStatusCodeExpelled
	default:
		key = entity.EducationStatusCodeEntrance
	}
	return "master.educationStatus." + util.ToCamel(key)
}

func degreeTypeKey(degreeTypeID uint64) string {
	var key string
	switch degreeTypeID {
	case entity.DegreeTypeHighSchool:
		key = entity.DegreeTypeCodeHighSchool
	case entity.DegreeTypeVocational:
		key = entity.DegreeTypeCodeVocational
	case entity.DegreeTypeJuniorCollege:
		key = entity.DegreeTypeCodeJuniorCollege
	case entity.DegreeTypeBachelor:
		key = entity.DegreeTypeCodeBachelor
	case entity.DegreeTypeMaster:
		key = entity.DegreeTypeCodeMaster
	case entity.DegreeTypeDoctor:
		key = entity.DegreeTypeCodeDoctor
	default:
		key = entity.DegreeTypeCodeOther
	}
	return "master.degreeType." + util.ToCamel(key)
}

// ToUserEducationRes は エンティティをレスポンスDTOに変換する
func ToUserEducationRes(e *entity.UserEducation) response.UserEducationResponse {
	return response.UserEducationResponse{
		ID:                e.ID,
		UserID:            e.UserID,
		InstitutionName:   e.InstitutionName,
		FacultyName:       e.FacultyName,
		DepartmentName:    e.DepartmentName,
		DegreeTypeID:      e.DegreeTypeID,
		DegreeType:        degreeTypeKey(e.DegreeTypeID),
		EducationStatusID: e.EducationStatusID,
		EducationStatus:   educationStatusKey(e.EducationStatusID),
		EventDate:         e.EventDate.Format("2006-01"),
		Description:       e.Description,
		SortOrder:         e.SortOrder,
		IsPublic:          e.IsPublic,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

// ToUserEducationResponses は 学歴エンティティ配列をレスポンスDTOに変換する
func ToUserEducationResponses(es []*entity.UserEducation) []response.UserEducationResponse {
	out := make([]response.UserEducationResponse, 0, len(es))
	for _, e := range es {
		out = append(out, ToUserEducationRes(e))
	}
	return out
}

// assembleHasUserEducation は存在有無のプリミティブを出力 DTO に変換します。
func assembleHasUserEducation(exists bool) HasEducationOutput {
	return HasEducationOutput{
		Exists: exists,
	}
}
