// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import (
	"context"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
)

// ListCountry は、国の選択肢一覧を取得するユースケースです。
func (uc *Interactor) ListCountry(ctx context.Context, in ListCountryInput) (ListCountryOutput, error) {
	var list []entity.Country
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		list, err = uc.cnRepo.List(txCtx, in.OnlySupported)
		return err
	}); err != nil {
		return ListCountryOutput{}, err
	}
	out := assembleCountry(list, func(code, fallback string) string {
		// master/country.yaml → master.country.<code>
		key := vo.NewKey("master.country." + code)
		label, ok := uc.translator.Translate(vo.NewLocale(in.Locale), key)
		if !ok {
			if fallback != "" {
				return fallback
			}
			return code
		}
		return label
	})
	return out, nil
}
