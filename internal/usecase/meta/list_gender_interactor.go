// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import (
	"context"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
)

// ListGender は、性別の選択肢一覧を取得するユースケースです。
func (uc *Interactor) ListGender(ctx context.Context, in ListGenderInput) (ListGenderOutput, error) {
	var list []entity.Gender
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		list, err = uc.gnRepo.ListActive(txCtx)
		return err
	}); err != nil {
		return ListGenderOutput{}, err
	}
	out := assembleGender(list, func(code, fallback string) string {
		// master/gender.yaml → master.gender.<code>
		key := vo.NewKey("master.gender." + code)

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
