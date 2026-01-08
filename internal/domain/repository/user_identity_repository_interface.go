// Package repository は 認証プロバイダ情報のインフラ層抽象
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// ListUserIdentitySpec は、一覧取得（リスト）系ユースケースで使用される共通の検索条件インターフェースです。
// ページング、ソート、キーワード検索などの仕様を統一的に扱うために利用されます。
type ListUserIdentitySpec interface {
	BaseListSpec

	// Q は、フリーテキスト検索用のクエリ文字列を返します。
	// 検索語が指定されていない場合は nil を返します。
	Q() *string
}

// UserIdentityRepository はログインプロバイダの情報を扱うリポジトリインターフェースです。
type UserIdentityRepository interface {
	// FindByProviderUID は Provider + ProviderUserID で一意の認証情報を取得します。
	FindByProviderUID(ctx context.Context, provider, providerUID string) (*entity.UserIdentity, error)

	// Upsert は既存なら更新、存在しなければ新規作成します。
	Upsert(ctx context.Context, ent *entity.UserIdentity) (*entity.UserIdentity, error)

	// ListByUserID は指定ユーザー配下の認証情報一覧を取得します。
	ListByUserID(ctx context.Context, userID uint64) ([]*entity.UserIdentity, error)

	// ListByUserIDWithSpec は指定ユーザー配下の認証情報一覧を取得します。
	ListByUserIDWithSpec(ctx context.Context, userID uint64, spec ListUserIdentitySpec) ([]*entity.UserIdentity, int64, error)
}
