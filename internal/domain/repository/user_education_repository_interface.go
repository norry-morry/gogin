// Package repository は、ドメイン層における永続化のためのリポジトリインターフェースを定義します。
// このファイルでは、ユーザー学歴を扱うリポジトリの契約を定義します。
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// UserEducationRepository は ユーザー学歴を扱うリポジトリインターフェースです
// このテーブルは原則として、必ずuser_idをwhere句に使う前提で想定しています。(暫定)
type UserEducationRepository interface {
	// HasAnyEducation は 指定ユーザーが学歴の登録をしているか判定します
	HasAnyEducation(ctx context.Context, userID uint64) (bool, error)

	// FindByID は 学歴を1件取得します
	FindByID(ctx context.Context, ID uint64) (*entity.UserEducation, error)

	// FindByUserID は 指定ユーザーの学歴を1件取得します
	FindByUserID(ctx context.Context, userID uint64) (*entity.UserEducation, error)

	// ListByUserID は 指定ユーザーの学歴リストを取得します
	ListByUserID(ctx context.Context, userID uint64) ([]*entity.UserEducation, error)

	// CountByUserID は 指定ユーザーの学歴登録件数を取得します(SortOrder計算用)
	CountByUserID(ctx context.Context, userID uint64) (int, error)

	// Create は 学歴を作成します
	Create(ctx context.Context, edu *entity.UserEducation) (*entity.UserEducation, error)

	// Update は 学歴を更新します
	Update(ctx context.Context, edu *entity.UserEducation) (*entity.UserEducation, error)

	// Delete は 学歴を物削します
	Delete(ctx context.Context, ID uint64) error

	// UpdateOrder は 学歴を並び替えします
	UpdateOrder(ctx context.Context, userID uint64, educationIDs []uint64) error
}
