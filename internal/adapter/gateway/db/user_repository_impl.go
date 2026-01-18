// Package db は GORM を使用したユーザーリポジトリ実装を提供します。
package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

// GormUserRepository は UserRepository の GORM 実装です。
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository は UserRepository インターフェースの GORM 実装を返します。
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{db: db}
}

//func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
//	var users []*entity.User
//	if err := r.db.WithContext(ctx).
//		Where("deleted_at IS NULL").
//		Find(&users).Error; err != nil {
//		return nil, err
//	}
//	return users, nil
//}

// FindByID は ID でユーザーを取得します。
func (r *GormUserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var u entity.User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUID は UID でユーザーを取得します。
// 見つからない場合は (nil, nil) を返します。
func (r *GormUserRepository) FindByUID(ctx context.Context, uid string) (*entity.User, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var u entity.User
	err := db.Where("uid = ?", uid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail はメールアドレスでユーザーを取得します。
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var u entity.User
	err := db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create はユーザーを作成します。
func (r *GormUserRepository) Create(ctx context.Context, u *entity.User) (uint64, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	if err := db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

// UpdateProfileOnLogin はログイン時のプロフィール項目を更新します。
func (r *GormUserRepository) UpdateProfileOnLogin(ctx context.Context, u *entity.User) error {
	db := db2.FromCtxOrDB(ctx, r.db)
	updates := map[string]any{
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"display_name":   u.DisplayName,
		"photo_url":      u.PhotoURL,
		"disabled":       u.Disabled,
		"last_login_at":  u.LastLoginAt,
	}
	return db.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Updates(updates).Error
}

//func (r *UserRepository) Update(ctx context.Context, u *entity.User) error {
//	return r.db.Save(u).Error
//}

//func (r *UserRepository) Delete(id uint) error {
//	return r.db.Delete(&entity.User{}, id).Error
//}

//func (r *UserRepository) FindByIdentity(ctx context.Context, provider, providerUserID string) (*entity.User, error) {
//	var u entity.User
//	// users と auth_identities を結合して取得（退会済みは除外）
//	err := r.db.WithContext(ctx).
//		Table("users").
//		Joins("INNER JOIN auth_identities ai ON ai.user_id = users.id").
//		Where("ai.provider = ? AND ai.provider_user_id = ?", provider, providerUserID).
//		Where("users.deleted_at IS NULL").
//		Limit(1).
//		Take(&u).Error
//
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, gorm.ErrRecordNotFound
//		}
//		return nil, err
//	}
//	return &u, nil
//}
