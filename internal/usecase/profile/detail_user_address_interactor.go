package profile

import (
	"context"

	"resume/internal/domain/entity"
	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// DetailUserAddress は 認証済みユーザーの住所情報を返します
func (uc *Interactor) DetailUserAddress(
	ctx context.Context,
	in DetailUserAddressInput,
) (DetailUserAddressOutput, error) {
	var addr *entity.UserAddress
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		a, err := uc.uaRepo.FindByID(txCtx, in.AddressID)
		if err != nil {
			return err
		}

		if a.UserID != in.UserID {
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
		addr = a
		return nil
	}); err != nil {
		return DetailUserAddressOutput{}, err
	}

	return DetailUserAddressOutput{
		UserAddress: addr,
	}, nil
}
