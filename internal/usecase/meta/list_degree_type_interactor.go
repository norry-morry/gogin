// Package meta は 学位種別のメタ情報を扱うユースケースを提供します
package meta

import (
	"context"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
)

// ListDegreeType は 学位種別のメタ情報を扱うユースケースです
func (uc *Interactor) ListDegreeType(
	ctx context.Context,
	in ListDegreeTypeInput,
) (ListDegreeTypeOutput, error) {
	var list []entity.DegreeType

	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		list, err = uc.dtRepo.ListIsActive(txCtx)
		return err
	}); err != nil {
		return ListDegreeTypeOutput{}, err
	}

	out := assembleDegreeType(list, func(code, fallback string) string {
		key := vo.NewKey("master.degree_type." + code)

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
