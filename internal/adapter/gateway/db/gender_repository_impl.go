// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

type genderRepository struct {
	db *gorm.DB
}

// NewGenderRepository は GenderRepository インターフェースのGORM実装を返します
func NewGenderRepository(db *gorm.DB) repository.GenderRepository {
	return &genderRepository{db: db}
}

// ListActive は 有効な性別一覧を返します
func (r *genderRepository) ListActive(ctx context.Context) ([]entity.Gender, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var list []entity.Gender
	if err := db.
		Where("is_active = ?", true).
		Order("sort_order, code").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
