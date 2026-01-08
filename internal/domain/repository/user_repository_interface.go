// Package repository は フリーランサー情報のインフラ層抽象
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// UserRepository はユーザー情報へのリポジトリ操作を定義するインターフェースです。
type UserRepository interface {
	// FindByID はユーザーをidで取得します
	FindByID(ctx context.Context, id uint64) (*entity.User, error)

	// FindByUID はユーザーをuidで取得します
	FindByUID(ctx context.Context, uid string) (*entity.User, error)

	// FindByEmail はユーザーをメールアドレスで取得します
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Create はユーザーを作成します
	Create(ctx context.Context, u *entity.User) (uint64, error)

	// UpdateProfileOnLogin はユーザーを更新します
	UpdateProfileOnLogin(ctx context.Context, u *entity.User) error

	// FindAll(ctx context.Context) ([]*entity.User, error)
	// FindByIdentity(ctx context.Context, provider, providerUserID string) (*entity.User, error)
	// Update(ctx context.Context, u *entity.User) error
	// Delete(id uint) error
}
