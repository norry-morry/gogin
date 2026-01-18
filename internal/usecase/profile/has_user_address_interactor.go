package profile

import "context"

// HasUserAddress は HasUserAddress ユースケースの実装本体です。
// Interactor.HasUserAddress から呼ばれる private 実装にしています。
func (uc *Interactor) HasUserAddress(ctx context.Context, in UserInput) (HasUserAddressOutput, error) {
	if in.UserID == 0 {
		return HasUserAddressOutput{}, errUnauthorized()
	}

	var exists bool
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		exists, err = uc.uaRepo.HasAnyAddress(txCtx, in.UserID)
		return err
	}); err != nil {
		return HasUserAddressOutput{}, errCheckAddressFailed(err)
	}
	return assembleHasUserAddressOutput(exists), nil
}
