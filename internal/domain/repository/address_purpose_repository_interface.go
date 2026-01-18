package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// AddressPurposeRepository は住所の目的を扱うリポジトリインターフェースです。
type AddressPurposeRepository interface {
	// GetAllActive は 有効な目的一覧を取得します
	GetAllActive(ctx context.Context) ([]entity.AddressPurpose, error)
}
