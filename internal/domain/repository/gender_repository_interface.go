package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// GenderRepository は 性別を扱うリポジトリインターフェースです。
type GenderRepository interface {
	ListActive(ctx context.Context) ([]entity.Gender, error)
}
