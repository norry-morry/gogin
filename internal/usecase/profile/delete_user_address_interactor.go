package profile

import (
	"context"

	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// DeleteUserAddress は 認証済みユーザーの住所情報を削除します
func (uc *Interactor) DeleteUserAddress(
	ctx context.Context,
	in DeleteUserAddressInput,
) (DeleteUserAddressOutput, error) {
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		addr, err := uc.uaRepo.FindByID(txCtx, in.AddressID)
		if err != nil {
			return err
		}

		if addr.UserID != in.UserID {
			ve := valerr.New()
			ve.Add(
				"domain.user_address.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.userAddress.id",
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}

		return uc.uaRepo.SoftDelete(txCtx, addr.ID)
	}); err != nil {
		return DeleteUserAddressOutput{}, err
	}

	return DeleteUserAddressOutput{}, nil
}
