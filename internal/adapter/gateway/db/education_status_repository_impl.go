// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

type educationStatusRepository struct {
	db *gorm.DB
}

// NewEducationStatusRepository は EducationStatusRepository interfaceのGORM実装を返します
func NewEducationStatusRepository(db *gorm.DB) repository.EducationStatusRepository {
	return &educationStatusRepository{db: db}
}

// ListIsActive は 有効な学歴状態一覧を取得染ます
func (r *educationStatusRepository) ListIsActive(
	ctx context.Context,
) ([]entity.EducationStatus, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	var list []entity.EducationStatus
	if err := db.
		Where("is_active = ?", true).
		Order("sort_order, code").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
