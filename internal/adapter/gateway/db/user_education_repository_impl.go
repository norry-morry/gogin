// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
	"resume/internal/infra/db/dberr"
)

type userEducationRepository struct {
	db *gorm.DB
}

// NewUserEducationRepository は UserEducationRepository interfaceのGORM実装を返します
func NewUserEducationRepository(db *gorm.DB) repository.UserEducationRepository {
	return &userEducationRepository{db: db}
}

// HasEducation は 指定ユーザーが学歴の登録をしているか判定します

// FindByID は 学歴を1件取得します
func (r *userEducationRepository) FindByID(
	ctx context.Context,
	ID uint64,
) (*entity.UserEducation, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var ue entity.UserEducation

	err := db.
		Where("id = ?", ID).
		First(&ue).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ue, nil
}

// FindByUserID は 指定ユーザーの学歴を1件取得します
func (r *userEducationRepository) FindByUserID(
	ctx context.Context,
	userID uint64,
) (*entity.UserEducation, error) {
	db := db2.
		FromCtxOrDB(ctx, r.db).
		Model(&entity.UserEducation{})

	var ue entity.UserEducation
	err := db.
		Where("user_id = ?", userID).
		Limit(1).
		Take(&ue).
		Error

	switch {
	case err == nil:
		return &ue, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// ListByUserID は 指定ユーザーの学歴リストを取得します
func (r *userEducationRepository) ListByUserID(
	ctx context.Context,
	userID uint64,
) ([]*entity.UserEducation, error) {
	db := db2.FromCtxOrDB(ctx, r.db).Model(&entity.UserEducation{})
	var rows []*entity.UserEducation
	db = db.
		Where("user_id = ?", userID). // AND is_public = true
		Order("sort_order ASC, id ASC")
	if err := db.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CountByUserID は 指定ユーザーの学歴登録件数を取得します(SortOrder計算用)
func (r *userEducationRepository) CountByUserID(
	ctx context.Context,
	userID uint64,
) (int, error) {
	db := db2.
		FromCtxOrDB(ctx, r.db).
		Model(&entity.UserEducation{})
	var count int64

	if err := db.Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// Create は 学歴を作成します
func (r *userEducationRepository) Create(
	ctx context.Context,
	edu *entity.UserEducation,
) (*entity.UserEducation, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	if err := db.Create(edu).Error; err != nil {
		return nil, dberr.Map(err)
	}
	return edu, nil
}

// Update は 学歴を更新します
func (r *userEducationRepository) Update(
	ctx context.Context,
	edu *entity.UserEducation,
) (*entity.UserEducation, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	if err := db.
		Model(&entity.UserEducation{}).
		Where("id = ? AND user_id = ?", edu.ID, edu.UserID).
		Updates(edu).
		Error; err != nil {
		return nil, dberr.Map(err)
	}
	return edu, nil
}

// Delete は 学歴を物削します
func (r *userEducationRepository) Delete(
	ctx context.Context,
	ID uint64,
) error {
	db := db2.FromCtxOrDB(ctx, r.db)

	return db.
		Delete(&entity.UserEducation{}, ID).
		Error
}

// UpdateOrder は 学歴を並び替えします
func (r *userEducationRepository) UpdateOrder(
	ctx context.Context,
	userID uint64,
	educationIDs []uint64,
) error {
	db := db2.FromCtxOrDB(ctx, r.db)

	// 取得したDBを使ってトランザクションを開始
	// ループで毎回コミットとすると途中で転けたときに回復できなく成るので
	return db.Transaction(func(tx *gorm.DB) error {
		for i, id := range educationIDs {
			sortOrder := i + 1
			if err := tx.
				Model(&entity.UserEducation{}).
				Where("id = ?", id).
				Where("user_id = ?", userID).
				Update("sort_order", sortOrder).
				Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// HasAnyEducation は 指定ユーザーが学歴の登録をしているか判定します
func (r *userEducationRepository) HasAnyEducation(
	ctx context.Context,
	userID uint64,
) (bool, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var ue entity.UserEducation
	err := db.
		Select("id").
		Where("user_id = ?", userID).
		Limit(1).
		Take(&ue).
		Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}
