package education

import "context"

// HasUserEducation は 認証済みユーザーの学歴の有無を返します
func (uc *Interactor) HasUserEducation(
	ctx context.Context,
	in HasEducationInput,
) (HasEducationOutput, error) {
	if in.UserID == 0 {
		return HasEducationOutput{}, errUnauthorized()
	}

	var exists bool
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		exists, err = uc.ueRepo.HasAnyEducation(txCtx, in.UserID)
		return err
	}); err != nil {
		return HasEducationOutput{}, err
	}
	return assembleHasUserEducation(exists), nil
}
