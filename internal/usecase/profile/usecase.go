// Package profile は ログイン後の自身の情報を取得する実装
package profile

import "context"

// Usecase は 認証ユーザーユースケースの操作を定義します
// 現在想定しているものとしては、
type Usecase interface {
	// GetUser は usersから当該の1行取得
	GetUser(ctx context.Context, in UserInput) (UserOutput, error)

	// ListUserIdentity は user_identitiesから当該の複数行
	ListUserIdentity(ctx context.Context, in UserIdentityInput) (*ListUserIdentityOutput, error)

	// GetUserProfile は 履歴書記載の情報の一部・氏名とか生年月日とかを取得
	GetUserProfile(ctx context.Context, in UserInput) (UserProfileOutput, error)

	// PatchUserProfile は 履歴書記載の情報の一部・氏名とか生年月日とかを更新
	PatchUserProfile(ctx context.Context, in PatchUserProfileInput) error

	// HasUserProfile は 履歴書記載の情報の一部・氏名とか生年月日が登録済みか否かを取得
	HasUserProfile(ctx context.Context, in UserInput) (HasUserProfileOutput, error)

	// HasUserAddress は ご住所を登録済みか否かを取得
	HasUserAddress(ctx context.Context, in UserInput) (HasUserAddressOutput, error)

	// ListUserAddress は ご住所一覧を取得
	ListUserAddress(ctx context.Context, in ListUserAddressInput) (ListUserAddressOutput, error)

	// DetailUserAddress は ご住所を取得
	DetailUserAddress(ctx context.Context, in DetailUserAddressInput) (DetailUserAddressOutput, error)

	// CreateUserAddress は ご住所を作成
	CreateUserAddress(ctx context.Context, in CreateUserAddressInput) (CreateUserAddressOutput, error)

	// UpdateUserAddress は ご住所を作成
	UpdateUserAddress(ctx context.Context, in UpdateUserAddressInput) (*UpdateUserAddressOutput, error)

	// DeleteUserAddress は ご住所を作成
	DeleteUserAddress(ctx context.Context, in DeleteUserAddressInput) (DeleteUserAddressOutput, error)
}
