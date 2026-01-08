package user

import (
	"context"

	"resume/internal/domain/repository"
	"resume/internal/shared/apperr"
	ucshared "resume/internal/usecase/shared"
)

type interactor struct {
	repo repository.UserRepository
}

// NewUsecase はユーザユースケースの実装を生成します。
func NewUsecase(repo repository.UserRepository) Usecase {
	return &interactor{repo: repo}
}

func (i *interactor) GetByID(ctx context.Context, in GetUserInput) (*UserOutput, error) {
	u, err := i.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, ucshared.MapInfraError(err) // 既存の共通マッパに委譲
	}
	if u == nil {
		return nil, apperr.New(apperr.CodeNotFound, "user not found", map[string]any{"id": in.ID})
	}
	return &UserOutput{
		ID:            u.ID,
		UID:           u.UID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		DisplayName:   u.DisplayName,
		PhotoURL:      u.PhotoURL,
		Disabled:      u.Disabled,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     u.DeletedAt,
	}, nil
}
