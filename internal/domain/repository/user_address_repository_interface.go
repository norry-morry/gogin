// Package repository は フリーランサー住所のインフラ層抽象
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// ListUserAddressSpec は、一覧取得（リスト）系ユースケースで使用される共通の検索条件インターフェースです。
// ページング、ソート、キーワード検索などの仕様を統一的に扱うために利用されます。
type ListUserAddressSpec interface {
	BaseListSpec
	PurposeIDVal() *uint64
}

// UserAddressRepository はフリーランサー住所情報を扱うリポジトリインターフェースです。
type UserAddressRepository interface {

	// FindByID は 指定IDのユーザー住所を取得します
	FindByID(ctx context.Context, ID uint64) (*entity.UserAddress, error)

	// HasAnyAddress は指定ユーザーが住所を登録しているかを判定します。
	HasAnyAddress(ctx context.Context, userID uint64) (bool, error)

	// ListUserAddresses はユーザー住所一覧を取得します。
	ListUserAddresses(ctx context.Context, userID uint64, onlyPrimary bool, purposeCode *string) ([]entity.UserAddress, error)

	// Create は ユーザー住所を作成します
	Create(ctx context.Context, a *entity.UserAddress) (*entity.UserAddress, error)

	// ExistsByUserIDAndPurposeID は 目的の重複を確認します
	ExistsByUserIDAndPurposeID(ctx context.Context, userID uint64, purposeID uint64) (bool, error)

	// ExistsByUserIDAndPurposeIDExceptID は、指定した住所ID以外で (user_id, purpose_id) が存在するかを確認します。
	ExistsByUserIDAndPurposeIDExceptID(
		ctx context.Context,
		userID uint64,
		purposeID uint64,
		exceptAddressID uint64,
	) (bool, error)

	// ListByUserIDWithSpec は 指定ユーザー配下の住所一覧を取得します
	ListByUserIDWithSpec(ctx context.Context, userID uint64, spec ListUserAddressSpec) ([]*entity.UserAddress, int64, error)

	// Update は ユーザー住所を更新します
	// a.ID, a.UserID などをキーに UPDATE する実装を想定。
	Update(ctx context.Context, a *entity.UserAddress) (*entity.UserAddress, error)

	// SoftDelete は ユーザー住所を論削します
	SoftDelete(ctx context.Context, id uint64) error
}
