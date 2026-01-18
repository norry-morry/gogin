// Package meta は 学歴状態のメタ情報を扱うユースケースを提供します
package meta

import (
	"context"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
)

// ListEducationStatus は 学歴状態のメタ情報を扱うユースケースです
func (uc *Interactor) ListEducationStatus(
	ctx context.Context,
	in ListEducationStatusInput,
) (ListEducationStatusOutput, error) {
	var list []entity.EducationStatus

	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		var err error
		list, err = uc.esRepo.ListIsActive(txCtx)
		return err
	}); err != nil {
		return ListEducationStatusOutput{}, err
	}

	out := assembleEducationStatus(list, func(code, fallback string) string {
		key := vo.NewKey("master.education_status." + code)

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
