// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

type degreeTypeRepository struct {
	db *gorm.DB
}

// NewDegreeTypeRepository は DegreeTypeRepository インターフェースのGORM実装を返します
func NewDegreeTypeRepository(db *gorm.DB) repository.DegreeTypeRepository {
	return &degreeTypeRepository{
		db: db,
	}
}

// ListIsActive は 有効な学位種別一覧を返します
func (r *degreeTypeRepository) ListIsActive(
	ctx context.Context,
) ([]entity.DegreeType, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	var list []entity.DegreeType
	if err := db.
		Where("is_active = ?", true).
		Order("sort_order, code").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
