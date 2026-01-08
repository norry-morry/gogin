// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import (
	"context"

	"resume/internal/domain/entity"
)

// CreateUserEducation は、指定されたユーザーIDに紐づく 学歴情報 を登録します。
// トランザクション内で UserEducation を参照し、並び順 を返します。
func (uc *Interactor) CreateUserEducation(
	ctx context.Context,
	in CreateInput,
) (CreateOutput, error) {
	edu := &entity.UserEducation{
		UserID:            in.UserID,
		InstitutionName:   in.InstitutionName,
		FacultyName:       in.FacultyName,
		DepartmentName:    in.DepartmentName,
		DegreeTypeID:      in.DegreeTypeID,
		EducationStatusID: in.EducationStatusID,
		EventDate:         in.EventDate,
		Description:       in.Description,
		IsPublic:          in.IsPublic,
	}
	var saved *entity.UserEducation
	err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		count, err := uc.ueRepo.CountByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}
		edu.SortOrder = count

		res, err := uc.ueRepo.Create(txCtx, edu)
		if err != nil {
			return err
		}
		saved = res
		return nil
	})
	if err != nil {
		return CreateOutput{}, err
	}
	return CreateOutput{
		Education: saved,
	}, err
}
