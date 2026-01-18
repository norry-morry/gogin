// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
	"resume/internal/infra/db/dberr"
)

type addressRepository struct {
	db *gorm.DB
}

// NewUserAddressRepository は UserAddressRepository インターフェースのGORM実装を返します。
func NewUserAddressRepository(db *gorm.DB) repository.UserAddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) FindByID(
	ctx context.Context,
	ID uint64,
) (*entity.UserAddress, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var ua entity.UserAddress

	err := db.
		Where("id = ?", ID).
		First(&ua).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ua, nil
}

func (r *addressRepository) HasAnyAddress(ctx context.Context, userID uint64) (bool, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var ua entity.UserAddress
	err := db.Select("id").Where("user_id = ?", userID).Limit(1).Take(&ua).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (r *addressRepository) ListUserAddresses(ctx context.Context, userID uint64, onlyPrimary bool, purposeCode *string) ([]entity.UserAddress, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	tx := db.Model(&entity.UserAddress{}).
		Preload("Purpose").
		Where("user_id = ?", userID)

	if onlyPrimary {
		tx = tx.Where("is_primary = ?", true)
	}
	if purposeCode != nil && *purposeCode != "" {
		// テーブル名は実テーブル/ネーミング戦略に合わせて要確認
		tx = tx.Joins("JOIN address_purposes ap ON ap.id = user_addresses.purpose_id").
			Where("ap.code = ?", *purposeCode)
		// 将来的な重複に備えて DISTINCT を付けておくと堅い
		tx = tx.Distinct("user_addresses.id")
	}

	var list []entity.UserAddress
	if err := tx.Order("user_addresses.id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *addressRepository) Create(ctx context.Context, a *entity.UserAddress) (*entity.UserAddress, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	if err := db.Create(a).Error; err != nil {
		return nil, dberr.Map(err)
	}
	return a, nil
}

func (r *addressRepository) ExistsByUserIDAndPurposeID(ctx context.Context, userID uint64, purposeID uint64) (bool, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	var count int64

	err := db.Model(&entity.UserAddress{}).
		Where("user_id = ? AND purpose_id = ?", userID, purposeID).
		Limit(1).
		Count(&count).Error

	if err != nil {
		return false, dberr.Map(err)
	}

	return count > 0, nil
}

func (r *addressRepository) ListByUserIDWithSpec(
	ctx context.Context,
	userID uint64,
	spec repository.ListUserAddressSpec,
) ([]*entity.UserAddress, int64, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserAddress{})
	db = db.Where("user_id = ?", userID)
	if purposeID := spec.PurposeIDVal(); purposeID != nil && *purposeID > 0 {
		db = db.Where("purpose_id = ?", purposeID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	qdb := db.Order(fmt.Sprintf("%s %s", spec.SortCol(), spec.Order()))
	if spec.Limit() > 0 {
		qdb = qdb.Limit(spec.Limit()).Offset(spec.Offset())
	}
	var rows []*entity.UserAddress
	if err := qdb.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *addressRepository) Update(
	ctx context.Context,
	a *entity.UserAddress,
) (*entity.UserAddress, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	//if err := db.Save(a).Error; err != nil {
	//	return nil, dberr.Map(err)
	//}
	if err := db.
		Model(&entity.UserAddress{}).
		Where("id = ? AND user_id = ?", a.ID, a.UserID).
		Updates(a).Error; err != nil {
		return nil, dberr.Map(err)
	}
	return a, nil
}

func (r *addressRepository) ExistsByUserIDAndPurposeIDExceptID(
	ctx context.Context,
	userID uint64,
	purposeID uint64,
	exceptAddressID uint64,
) (bool, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	var count int64
	if err := db.
		Model(&entity.UserAddress{}).
		Where("user_id = ? AND purpose_id = ? AND id <> ?", userID, purposeID, exceptAddressID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// SoftDelete は ユーザー住所を論削します
func (r *addressRepository) SoftDelete(
	ctx context.Context,
	id uint64,
) error {
	db := db2.FromCtxOrDB(ctx, r.db)

	return db.
		Where("id = ?", id).
		Delete(&entity.UserAddress{}).
		Error
}
