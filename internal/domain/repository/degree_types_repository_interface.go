// Package repository は、ドメイン層における永続化のためのリポジトリインターフェースを定義します。
// このファイルでは、学位種別を扱うリポジトリの契約を定義します。
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// DegreeTypeRepository は 学位種別を扱うリポジトリインターフェースです
type DegreeTypeRepository interface {
	// ListIsActive は 学位状態一覧を取得します
	ListIsActive(ctx context.Context) ([]entity.DegreeType, error)
}
