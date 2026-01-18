package education

import (
	"context"

	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// ReorderEducation は 認証済みユーザーの学歴 を 並び順を変更 します
func (uc *Interactor) ReorderEducation(
	ctx context.Context,
	in ReorderInput,
) (ReorderOutput, error) {
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		// ここでEducationIDsが1個でもuserIDの物で無かったらエラーを返したい
		// 1. ユーザーが所有している学歴を全て取得
		ownedList, err := uc.ueRepo.ListByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}
		// 2. 所有しているIDのマップを作成(検索を高速化するため)
		ownedIDMap := make(map[uint64]bool)
		for _, edu := range ownedList {
			ownedIDMap[edu.ID] = true
		}
		// 3. リクエストされたIDが全てユーザーの物であるかチェック
		for _, reqID := range in.EducationIDs {
			if !ownedIDMap[reqID] {
				// 所有していないIDが含まれている場合はエラー
				ve := valerr.New()
				ve.Add(
					"domain.userEducation.id",
					"validation.rules.invalid",
					map[string]any{
						"field": "domain.userEducation.id",
						"value": reqID,
					},
				)
				return apperr.New(
					apperr.CodeUnprocessable,
					"The request contains semantically invalid data.",
					ve.ToDetails(),
				)
			}
		}
		// 4. チェックを通過したら並び替え実行
		return uc.ueRepo.UpdateOrder(txCtx, in.UserID, in.EducationIDs)
	}); err != nil {
		return ReorderOutput{}, err
	}
	return ReorderOutput{}, nil
}
