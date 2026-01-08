// Package repository は、ドメインモデルに対するリポジトリインターフェースを定義するパッケージです。
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// UserProfileRepository は、ユーザープロフィールへの永続化アクセスを提供するリポジトリインターフェースです。
type UserProfileRepository interface {
	// FindByUserID は、指定されたユーザーIDに対応するユーザープロフィールを1件取得します。
	// 対応するレコードが存在しない場合は、(*entity.UserProfile)(nil), nil を返します。
	FindByUserID(ctx context.Context, userID uint64) (*entity.UserProfile, error)
	Upsert(ctx context.Context, profile *entity.UserProfile) error
}
