// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"context"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

type countryRepository struct {
	db *gorm.DB
}

// NewCountryRepository は countryRepository インターフェースのGORM実装を返します
func NewCountryRepository(db *gorm.DB) repository.CountryRepository {
	return &countryRepository{db: db}
}

// List は 有効な国コード一覧を返します
func (r *countryRepository) List(
	ctx context.Context,
	onlySupported bool,
) ([]entity.Country, error) {
	db := db2.FromCtxOrDB(ctx, r.db)

	var list []entity.Country
	q := db.Model(&entity.Country{})
	if onlySupported {
		q = q.Where("is_supported = ?", true)
	}

	if err := q.Order("sort_order ASC, code ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
