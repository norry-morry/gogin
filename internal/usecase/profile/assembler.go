package profile

import (
	"time"

	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
)

func assembleUserOutput(user *entity.User) UserOutput {
	return UserOutput{User: user}
}

// ToUserIdentityRes は エンティティをレスポンスに変換するDTOです
func ToUserIdentityRes(e *entity.UserIdentity) response.UserIdentity {
	return response.UserIdentity{
		ID:                  e.ID,
		UserID:              e.UserID,
		Provider:            e.Provider,
		ProviderUserID:      e.ProviderUserID,
		ProviderDisplayName: e.ProviderDisplayName,
		EmailAtSignup:       e.EmailAtSignup,
		CreatedAt:           e.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserIdentityResponses は 複数のエンティティをレスポンスに変換するDTOです
func ToUserIdentityResponses(es []*entity.UserIdentity) []response.UserIdentity {
	out := make([]response.UserIdentity, 0, len(es))
	for _, e := range es {
		out = append(out, ToUserIdentityRes(e))
	}
	return out
}

// assembleHasAddressOutput は存在有無のプリミティブを出力 DTO に変換します。
func assembleHasUserAddressOutput(exists bool) HasUserAddressOutput {
	return HasUserAddressOutput{
		Exists: exists,
	}
}

func addressPurposeKey(purposeID uint64) string {
	switch purposeID {
	case entity.AddressPurposeIDHome:
		return "master.addressPurpose.home"
	case entity.AddressPurposeIDContact:
		return "master.addressPurpose.contact"
	case entity.AddressPurposeIDOffice:
		return "master.addressPurpose.office"
	case entity.AddressPurposeIDShipping:
		return "master.addressPurpose.shipping"
	case entity.AddressPurposeIDBilling:
		return "master.addressPurpose.billing"
	default:
		return "master.addressPurpose.other"
	}
}

// ToUserAddressRes は エンティティをレスポンス DTO に変換する
func ToUserAddressRes(e *entity.UserAddress) response.UserAddressResponse {
	return response.UserAddressResponse{
		ID:         e.ID,
		UserID:     e.UserID,
		PurposeID:  e.PurposeID,
		PurposeKey: addressPurposeKey(uint64(e.PurposeID)),
		IsPrimary:  e.IsPrimary,
		// Note: CountryCodeも多言語化対応した方が良いかも
		CountryCode:        e.CountryCode,
		AdministrativeArea: e.AdministrativeArea,
		Locality:           e.Locality,
		DependentLocality:  e.DependentLocality,
		PostalCode:         e.PostalCode,
		AddressLine1:       e.AddressLine1,
		AddressLine2:       e.AddressLine2,
		AddressLine3:       e.AddressLine3,
		Latitude:           e.Latitude,
		Longitude:          e.Longitude,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

// ToUserAddressResponses は ユーザー住所エンティティ配列をレスポンス DTO 配列に変換します。
func ToUserAddressResponses(es []*entity.UserAddress) []response.UserAddressResponse {
	out := make([]response.UserAddressResponse, 0, len(es))
	for _, e := range es {
		out = append(out, ToUserAddressRes(e))
	}
	return out
}
