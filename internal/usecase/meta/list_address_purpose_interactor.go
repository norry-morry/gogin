// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import (
	"context"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
)

// ListAddressPurpose は、住所用途の選択肢一覧を取得するユースケースです。
func (uc *Interactor) ListAddressPurpose(ctx context.Context, in ListAddressPurposeInput) (ListAddressPurposeOutput, error) {
	var list []entity.AddressPurpose
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		list, err = uc.apRepo.GetAllActive(txCtx)
		return err
	}); err != nil {
		return ListAddressPurposeOutput{}, err
	}

	out := assembleAddressPurpose(list, func(code, fallback string) string {
		// master/address_purpose.yaml → master.addressPurpose.<code>
		key := vo.NewKey("master.address_purpose." + code)

		label, ok := uc.translator.Translate(vo.NewLocale(in.Locale), key)
		if !ok {
			// fallback（YAML未定義なら DisplayName または code を返す）
			if fallback != "" {
				return fallback
			}
			return code
		}
		return label
	})
	return out, nil
}
