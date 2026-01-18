package profile

import (
	"context"
)

// ListUserIdentity は ListUserIdentity ユースケースの実装本体です。
// Interactor.ListUserIdentity から呼ばれる private 実装にしています。
func (uc *Interactor) ListUserIdentity(
	ctx context.Context,
	in UserIdentityInput,
) (*ListUserIdentityOutput, error) {
	entities, total, err := uc.uIRepo.ListByUserIDWithSpec(ctx, in.UserID, in)
	if err != nil {
		return nil, err
	}

	return &ListUserIdentityOutput{
		Items: ToUserIdentityResponses(entities),
		Total: total,
	}, nil
}
