// Package presenter は、ユースケースの出力データを HTTP レスポンス用 DTO に変換する責務を持ちます。
// 主にドメインエンティティ（entity）を response パッケージの構造体へ詰め替え、
// コントローラ層でそのまま JSON レスポンスとして返せる形式に整形します。
package presenter

import (
	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
)

// ToAddressResponse は、ドメインエンティティの UserAddress を
// HTTP レスポンス用の Address DTO に変換します。
// time.Time フィールドは JSON シリアライズ時に RFC3339 形式で出力されます。
func ToAddressResponse(a *entity.UserAddress) response.Address {
	return response.Address{
		ID:                 a.ID,
		PurposeID:          a.PurposeID,
		IsPrimary:          a.IsPrimary,
		CountryCode:        a.CountryCode,
		AdministrativeArea: a.AdministrativeArea,
		Locality:           a.Locality,
		DependentLocality:  a.DependentLocality,
		PostalCode:         a.PostalCode,
		AddressLine1:       a.AddressLine1,
		AddressLine2:       a.AddressLine2,
		AddressLine3:       a.AddressLine3,
		Latitude:           a.Latitude,
		Longitude:          a.Longitude,
		CreatedAt:          a.CreatedAt, // time.Time の JSON は RFC3339 系で出ます
		UpdatedAt:          a.UpdatedAt,
	}
}

// ToAddressListResponse は、UserAddress エンティティのスライスを
// Address DTO のスライスに変換します。
// 各要素は ToAddressResponse を用いて変換されます。
func ToAddressListResponse(list []*entity.UserAddress) []response.Address {
	out := make([]response.Address, 0, len(list))
	for _, a := range list {
		out = append(out, ToAddressResponse(a))
	}
	return out
}
