// Package profile は、ユーザープロフィールに関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package profile

import (
	"context"
	"errors"

	"resume/internal/domain/entity"
	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// CreateUserAddress は 認証済みユーザーの住所を新規登録します
// 入力値に基づきエンティティを生成し、永続化後に CreateUserAddressOutput を返します
// バリデーションエラー(業務ロジックエラー)などが発生した場合は適切なエラーを返します
func (uc *Interactor) CreateUserAddress(ctx context.Context, in CreateUserAddressInput) (CreateUserAddressOutput, error) {

	addr, err := entity.NewUserAddress(entity.UserAddressParam{
		UserID:             in.UserID,
		PurposeID:          in.PurposeID,
		IsPrimary:          in.IsPrimary,
		CountryCode:        in.CountryCode,
		AdministrativeArea: in.AdministrativeArea,
		Locality:           in.Locality,
		DependentLocality:  in.DependentLocality,
		PostalCode:         in.PostalCode,
		SortingCode:        in.SortingCode,
		AddressLine1:       in.AddressLine1,
		AddressLine2:       in.AddressLine2,
		AddressLine3:       in.AddressLine3,
		Latitude:           in.Latitude,
		Longitude:          in.Longitude,
	})
	if err != nil {
		var inv entity.InvalidAddressError
		if errors.As(err, &inv) {
			return CreateUserAddressOutput{}, apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				valerr.FromDomainAll(inv),
			)
		}
		return CreateUserAddressOutput{}, err
	}
	var saved *entity.UserAddress

	err = uc.tx.Do(ctx, func(txCtx context.Context) error {
		// 重複チェック(業務ルール)
		// 現住所の重複ができない様に確認して重複があればエラーを返す
		if in.PurposeID == entity.AddressPurposeIDHome {
			exists, err := uc.uaRepo.ExistsByUserIDAndPurposeID(txCtx, in.UserID, in.PurposeID)
			if err != nil {
				return err
			}
			if exists {
				ve := valerr.New()
				ve.Add(
					"domain.address.purposeId",
					"validation.rules.distinct",
					map[string]any{
						"field": "domain.address.purposeId.label",
					},
				)

				// 既存のバリデーションエラーと同じ Code/Message/Details で返す
				return apperr.New(
					apperr.CodeUnprocessable,
					"The request contains semantically invalid data.",
					ve.ToDetails(),
				)
			}
		}

		// INSERT実行
		s, err := uc.uaRepo.Create(txCtx, addr)
		if err != nil {
			return err
		}
		saved = s
		return nil
	})

	if err != nil {
		return CreateUserAddressOutput{}, err
	}

	return CreateUserAddressOutput{
		saved,
	}, nil
}
