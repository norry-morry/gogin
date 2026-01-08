// Package profile は、ユーザープロフィールに関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package profile

import (
	"context"

	"resume/internal/domain/entity"
	"resume/internal/shared/apperr"
	"resume/internal/shared/valerr"
)

// UpdateUserAddress は 認証済みユーザーの住所情報を更新します。
func (uc *Interactor) UpdateUserAddress(
	ctx context.Context,
	in UpdateUserAddressInput,
) (*UpdateUserAddressOutput, error) {
	var updated *entity.UserAddress

	err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		// ------------------------------------------------------
		// ① 対象レコード取得 → 自分の住所かチェック
		// ------------------------------------------------------
		current, err := uc.uaRepo.FindByID(txCtx, in.AddressID)
		if err != nil {
			return err
		}

		if current == nil {
			ve := valerr.New()
			ve.Add(
				"domain.user_address.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.user_address.id",
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}

		if current.UserID != in.UserID {
			ve := valerr.New()
			ve.Add(
				"domain.user_address.id",
				"validation.rules.invalid",
				map[string]any{
					"field": "domain.user_address.id",
				},
			)
			return apperr.New(
				apperr.CodeUnprocessable,
				"The request contains semantically invalid data.",
				ve.ToDetails(),
			)
		}
		// ------------------------------------------------------
		// ② 現住所(Home) の重複チェック
		//     ※ 自分自身のレコードは除外したいので ExceptID 版を使う
		// ------------------------------------------------------
		if in.PurposeID == entity.AddressPurposeIDHome {
			exists, err := uc.uaRepo.ExistsByUserIDAndPurposeIDExceptID(
				txCtx,
				in.UserID,
				in.PurposeID,
				in.AddressID,
			)
			if err != nil {
				return err
			}
			if exists {
				ve := valerr.New()
				ve.Add(
					"domain.user_address.purpose_id",
					"validation.rules.distinct",
					map[string]any{
						"field": "domain.user_address.purpose_id.label",
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

		// ------------------------------------------------------
		// ③ current に新しい値をマージして更新
		//    （新しく entity を作るより current を上書きするほうが自然）
		// ------------------------------------------------------
		current.PurposeID = in.PurposeID
		current.IsPrimary = in.IsPrimary
		current.CountryCode = in.CountryCode
		current.AdministrativeArea = in.AdministrativeArea
		current.Locality = in.Locality
		current.DependentLocality = in.DependentLocality
		current.PostalCode = in.PostalCode
		current.SortingCode = in.SortingCode
		current.AddressLine1 = in.AddressLine1
		current.AddressLine2 = in.AddressLine2
		current.AddressLine3 = in.AddressLine3
		current.Latitude = in.Latitude
		current.Longitude = in.Longitude

		updated, err = uc.uaRepo.Update(txCtx, current)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &UpdateUserAddressOutput{
		UserAddress: updated,
	}, nil
}
