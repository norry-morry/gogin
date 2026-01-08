// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

type addressPurposeRepository struct {
	db *gorm.DB
}

// NewAddressPurposeRepository は AddressPurposeRepository インターフェースのGORM実装を返します
func NewAddressPurposeRepository(db *gorm.DB) repository.AddressPurposeRepository {
	return &addressPurposeRepository{db: db}
}

// GetAllActive は 有効な目的一覧を返します
func (r *addressPurposeRepository) GetAllActive(ctx context.Context) ([]entity.AddressPurpose, error) {
	db := db2.FromCtxOrDB(ctx, r.db)
	var list []entity.AddressPurpose
	if err := db.
		Where("is_active = ?", true).
		Order("sort_order, code").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
