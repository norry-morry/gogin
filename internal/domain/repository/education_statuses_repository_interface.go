// Package repository は、ドメイン層における永続化のためのリポジトリインターフェースを定義します。
// このファイルでは、学歴状態を扱うリポジトリの契約を定義します。
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// EducationStatusRepository は 学歴状態を扱うリポジトリインターフェースです
type EducationStatusRepository interface {
	// ListIsActive は 有効な学歴状態一覧を取得染ます
	ListIsActive(ctx context.Context) ([]entity.EducationStatus, error)
}
