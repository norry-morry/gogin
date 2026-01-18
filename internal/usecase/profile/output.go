// Package profile は、ユーザープロフィールおよび認証関連のユースケース出力 DTO を定義します。
// ドメイン層のエンティティから変換されたデータを、Presenter 層や HTTP レスポンスへ渡すための構造体を提供します。
package profile

import (
	"time"

	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
)

// UserOutput は、ユーザーの認証情報（ログイン結果など）を表すユースケース層の出力 DTO です。
// entity.User を内包し、アプリケーション層から Presenter 層へ渡されます。
type UserOutput struct {
	User *entity.User
}

// ListUserIdentityOutput は、ユーザーに紐づく外部認証プロバイダ情報の一覧を表す出力 DTO です。
// ページング対応を想定し、items と total 件数を含みます。
type ListUserIdentityOutput struct {
	Items []response.UserIdentity `json:"items"`
	Total int64                   `json:"total"`
}

// UserProfileOutput は、ユーザーのプロフィール情報をユースケース層から Presenter 層へ渡すための出力 DTO です。
// 年齢（Age）はユースケース層で算出され、年齢帯（AgeGroup）は Presenter 層で導出されます。
// GenderLabelKey は、多言語化用の辞書キーを表します。
type UserProfileOutput struct {
	UserID         uint64     `json:"user_id"`
	FamilyName     string     `json:"family_name"`
	GivenName      string     `json:"given_name"`
	FamilyNameKana string     `json:"family_name_kana"`
	GivenNameKana  string     `json:"given_name_kana"`
	LegalName      string     `json:"legal_name"`
	LegalNameKana  string     `json:"legal_name_kana"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	Age            *int       `json:"age,omitempty"`       // ← usecase層で算出
	AgeGroup       *int       `json:"age_group,omitempty"` // ← presenter層で算出
	GenderID       *uint8     `json:"gender_id"`
	GenderLabelKey *string    `json:"gender_label_key"`
	Initial        *string    `json:"initial"`
}

// HasUserProfileOutput は、ユーザーのプロフィール情報有無をユースケース層から Presenter 層へ渡すための出力 DTO です。
type HasUserProfileOutput struct {
	Exists bool `json:"exists"`
}

// HasUserAddressOutput は、ユーザーの住所有無をユースケース層から Presenter 層へ渡すための出力 DTO です。
type HasUserAddressOutput struct {
	Exists bool `json:"exists"`
}

// CreateUserAddressOutput は ユーザーアドレスエンティティを言ったん暫定的に返します。
type CreateUserAddressOutput struct {
	UserAddress *entity.UserAddress `json:"user_address"`
}

// UpdateUserAddressOutput は ユーザーアドレスエンティティを言ったん暫定的に返します。
type UpdateUserAddressOutput struct {
	UserAddress *entity.UserAddress `json:"user_address"`
}

// ListUserAddressOutput は ユーザー住所一覧の出力 DTO を表します。
type ListUserAddressOutput struct {
	Items []response.UserAddressResponse `json:"items"`
	Total int64                          `json:"total"`
}

// DeleteUserAddressOutput は ユーザー住所削除の結果 DTO を表します。
type DeleteUserAddressOutput struct{}

// DetailUserAddressOutput は ユーザーアドレスエンティティを言ったん暫定的に返します。
type DetailUserAddressOutput struct {
	UserAddress *entity.UserAddress `json:"user_address"`
}
