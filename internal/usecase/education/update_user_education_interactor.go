// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import (
	"context"

	"resume/internal/domain/entity"
	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// UpdateUserEducation は、指定されたユーザーIDに紐づく 学歴情報 を更新します。
func (uc *Interactor) UpdateUserEducation(
	ctx context.Context,
	in UpdateInput,
) (UpdateOutput, error) {
	var updated *entity.UserEducation

	err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		current, err := uc.ueRepo.FindByID(txCtx, in.EducationID)
		if err != nil {
			return err
		}

		if current == nil {
			ve := valerr.New()
			ve.Add(
				"domain.user_education.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.user_education.id",
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}

		if current.UserID != in.UserID {
			ve := valerr.New()
			ve.Add(
				"domain.user_education.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.user_education.id",
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}

		current.InstitutionName = in.InstitutionName
		current.FacultyName = in.FacultyName
		current.DepartmentName = in.DepartmentName
		current.DegreeTypeID = in.DegreeTypeID
		current.EducationStatusID = in.EducationStatusID
		current.EventDate = in.EventDate
		current.Description = in.Description
		current.IsPublic = in.IsPublic

		res, err := uc.ueRepo.Update(txCtx, current)
		if err != nil {
			return err
		}
		updated = res
		return nil
	})
	if err != nil {
		return UpdateOutput{}, err
	}
	return UpdateOutput{
		Education: updated,
	}, nil
}
