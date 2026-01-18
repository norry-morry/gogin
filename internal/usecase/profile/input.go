// Package profile は 認容済みユーザーユースケースの I/O を定義します
package profile

import vo "resume/internal/usecase/valueobject"

// =======================
// Query inputs
// =======================

// UserInput は 認証済みユーザーのUserIdを入力します
type UserInput struct {
	UserID uint64
}

// UserIdentityInput は 認証済みユーザー検索のの入力項目を表す
type UserIdentityInput struct {
	UserID    uint64
	Page      int
	PerPage   int
	Sort      string
	SortOrder string
	Query     *string
}

// Offset は ページング機能のスタートを表す
func (in UserIdentityInput) Offset() int {
	if in.PerPage <= 0 || in.Page < 1 || in.PerPage == -1 {
		return 0
	}
	return (in.Page - 1) * in.PerPage
}

// Limit は ページング機能の表示件数を表す
func (in UserIdentityInput) Limit() int { return in.PerPage }

// SortCol は ページング機能のソート項目を表す
func (in UserIdentityInput) SortCol() string { return in.Sort }

// Order は ページング機能の昇降順を表す
func (in UserIdentityInput) Order() string { return in.SortOrder }

// Q は 検索波他メータを表す
func (in UserIdentityInput) Q() *string { return in.Query }

// PatchUserProfileInput は 認証済みユーザー個人情報の入力項目を表す
type PatchUserProfileInput struct {
	UserID         uint64
	FamilyName     string
	GivenName      string
	FamilyNameKana string
	GivenNameKana  string
	BirthDate      vo.PatchValue[string]
	GenderID       vo.PatchValue[uint8]
	Initial        vo.PatchValue[string]
}

// CreateUserAddressInput は ユーザーの住所登録入力項目を表す
type CreateUserAddressInput struct {
	UserID    uint64
	PurposeID uint64 // address_purposes.id（TINYINT UNSIGNED）
	IsPrimary bool   // 同一 (user_id, purpose_id) で true は最大 1 件

	// 位置情報（必須・任意の定義はマイグレーションに準拠）
	CountryCode        string  // CHAR(2) NOT NULL (ISO 3166-1 alpha-2)
	AdministrativeArea *string // VARCHAR(128) NULL
	Locality           *string // VARCHAR(128) NULL
	DependentLocality  *string // VARCHAR(128) NULL
	PostalCode         *string // VARCHAR(32)  NULL
	SortingCode        *string // VARCHAR(32)  NULL

	AddressLine1 string  // VARCHAR(160) NOT NULL
	AddressLine2 *string // VARCHAR(160) NULL
	AddressLine3 *string // VARCHAR(160) NULL

	Latitude  *float64 // DECIMAL(9,6)  NULL
	Longitude *float64 // DECIMAL(9,6)  NULL
}

// ListUserAddressInput は ユーザー住所一覧取得の入力項目を表します。
type ListUserAddressInput struct {
	UserID    uint64
	Page      int
	PerPage   int
	Sort      string
	SortOrder string

	PurposeID *uint64
	//Country   *string
	//City      *string
}

// Offset は ページング機能のスタートを表す
func (in ListUserAddressInput) Offset() int {
	if in.PerPage <= 0 || in.Page < 1 || in.PerPage == -1 {
		return 0
	}
	return (in.Page - 1) * in.PerPage
}

// Limit は ページング機能の表示件数を表す
func (in ListUserAddressInput) Limit() int { return in.PerPage }

// SortCol は ページング機能のソート項目を表す
func (in ListUserAddressInput) SortCol() string { return in.Sort }

// Order は ページング機能の昇降順を表す
func (in ListUserAddressInput) Order() string { return in.SortOrder }

// PurposeIDVal は 一覧検索仕様で利用する住所目的IDを返します。
func (in ListUserAddressInput) PurposeIDVal() *uint64 { return in.PurposeID }

// UpdateUserAddressInput は ユーザー住所更新の入力項目を表す
// CreateUserAddressInput を埋め込むことで、バリデーションやマッピングを共通化する
type UpdateUserAddressInput struct {
	AddressID uint64
	CreateUserAddressInput
}

// DeleteUserAddressInput は ユーザー住所の入力項目を表す
// 削除時に最低限のDTOでユーザー自身の住所かどうかを表す
type DeleteUserAddressInput struct {
	UserID    uint64
	AddressID uint64
}

// DetailUserAddressInput は ユーザー住所の入力項目を表す
// 詳細取得時に最低限のDTOでユーザー自身の住所かどうかを表す
type DetailUserAddressInput struct {
	UserID    uint64
	AddressID uint64
}
