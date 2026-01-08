package education

import (
	"context"

	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// DeleteUserEducation は 認証済みユーザーの学歴を削除します
// 削除完了後に並び順の再構築を行う
func (uc *Interactor) DeleteUserEducation(
	ctx context.Context,
	in DeleteInput,
) (DeleteOutput, error) {
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		edu, err := uc.ueRepo.FindByID(txCtx, in.EducationID)
		if err != nil {
			return err
		}
		if edu.UserID != in.UserID {
			ve := valerr.New()
			ve.Add(
				"domain.userEducation.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.userEducation.id",
					"value": in.EducationID,
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}
		// 1. Deleteの結果をハンドリング
		if err := uc.ueRepo.Delete(txCtx, in.EducationID); err != nil {
			return err
		}
		remainingList, err := uc.ueRepo.ListByUserID(txCtx, in.UserID)
		if err != nil {
			// ここは tx.Do の中なので error だけを返す
			return err
		}

		var newOrderIDs []uint64
		for _, edu := range remainingList {
			newOrderIDs = append(newOrderIDs, edu.ID)
		}

		// 2. UpdateOrderの結果を正しくハンドリング
		if len(newOrderIDs) > 0 {
			// if err := ...; err != nil の形に修正
			if err := uc.ueRepo.UpdateOrder(txCtx, in.UserID, newOrderIDs); err != nil {
				return err
			}
		}
		// 成功時は nil を返す
		return nil

	}); err != nil {
		return DeleteOutput{}, err
	}
	return DeleteOutput{}, nil
}
