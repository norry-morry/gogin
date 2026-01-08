// Package user はユーザー操作のユースケース層（ビジネスロジック）を定義します。
package user

import (
	"context"
)

// Usecase はユーザ領域のユースケースIFです。
type Usecase interface {
	// GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, int GetUserInput) (*UserOutput, error)
}
