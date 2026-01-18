package profile

import (
	"context"
	"errors"

	"resume/internal/domain/entity"
)

// GetUser は GetUser ユースケースの実装本体です。
// Interactor.GetUser から呼ばれる private 実装にしています。
func (uc *Interactor) GetUser(ctx context.Context, in UserInput) (UserOutput, error) {
	if in.UserID == 0 {
		return UserOutput{}, errUnauthorized()
	}
	var user *entity.User
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		u, err := uc.uRepo.FindByID(txCtx, in.UserID)
		if err != nil {
			// ここで repository.ErrNotFound をアプリ共通の NotFound に正規化できるとなお良い
			return err
		}
		if u == nil {
			// repository.ErrNotFound 等に差し替え可能
			return errors.New("not found")
		}
		user = u
		return nil
	}); err != nil {
		return UserOutput{}, err
	}

	return assembleUserOutput(user), nil
}
