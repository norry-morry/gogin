// Package db は GORM を使用した認証プロバイダ実装を提供します。
package db

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

// userIdentityRepository は UserIdentityRepository の GORM 実装です。
type userIdentityRepository struct {
	db *gorm.DB
}

// NewIdentityRepository は UserIdentityRepository インターフェースの
// GORM 実装を返します。
func NewIdentityRepository(db *gorm.DB) repository.UserIdentityRepository {
	return &userIdentityRepository{db: db}
}

// FindByProviderUID は provider + provider_user_id で認証IDを取得します。
// 見つからない場合は (nil, nil) を返します。
func (r *userIdentityRepository) FindByProviderUID(ctx context.Context, provider, providerUID string) (*entity.UserIdentity, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserIdentity{})
	//db = db.Debug()
	var ent entity.UserIdentity
	err := db.Where("provider = ? AND provider_user_id = ?", provider, providerUID).First(&ent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

// Upsert は provider + provider_user_id で一意な認証IDを upsert します。
// 既存があれば更新、無ければ作成します。
func (r *userIdentityRepository) Upsert(ctx context.Context, ent *entity.UserIdentity) (*entity.UserIdentity, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserIdentity{})
	//db = db.Debug()

	// GORM v2 の OnConflict 構文を利用（MySQL の INSERT ... ON DUPLICATE KEY UPDATE）
	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "provider"}, {Name: "provider_user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"user_id":               ent.UserID,
			"provider_display_name": ent.ProviderDisplayName,
			"email_at_signup":       ent.EmailAtSignup,
			// created_at は保持。必要なら updated_at があれば追記。
		}),
	}).Create(ent).Error; err != nil {
		return nil, err
	}
	return ent, nil
}

// ListByUserID は指定した user_id 配下の auth_identities 一覧を取得します。
func (r *userIdentityRepository) ListByUserID(ctx context.Context, userID uint64) ([]*entity.UserIdentity, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserIdentity{})
	//db = db.Debug()
	var rows []*entity.UserIdentity
	if err := db.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByUserIDWithSpec は指定した user_id 配下の auth_identities 一覧をページング要素含めて取得します。
func (r *userIdentityRepository) ListByUserIDWithSpec(
	ctx context.Context,
	userID uint64,
	spec repository.ListUserIdentitySpec,
) ([]*entity.UserIdentity, int64, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserIdentity{})
	//db = db.Debug()
	db = db.Where("user_id = ?", userID)

	if q := spec.Q(); q != nil && *q != "" {
		like := "%" + *q + "%"
		db = db.Where(
			"provider LIKE ? OR user_id LIKE ? OR email_at_signup LIKE ?",
			like,
			like,
			like,
		)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	qdb := db.Order(fmt.Sprintf("%s %s", spec.SortCol(), spec.Order()))
	if spec.Limit() > 0 { // -1 のときは全件
		qdb = qdb.Limit(spec.Limit()).Offset(spec.Offset())
	}

	var rows []*entity.UserIdentity
	if err := qdb.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
