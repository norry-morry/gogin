// Package profile は、ユーザープロフィールに関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package profile

import (
	"context"

	"resume/internal/shared/util"
)

// GetUserProfile は、指定されたユーザーIDに紐づくプロフィール情報を取得します。
// トランザクション内で UserProfileRepository を参照し、該当ユーザーが存在しない場合は空の DTO を返します。
// 取得したデータに対して年齢計算やラベルキー組み立てなどの補助的処理も行います。
func (uc *Interactor) GetUserProfile(ctx context.Context, in UserInput) (UserProfileOutput, error) {
	if in.UserID == 0 {
		return UserProfileOutput{}, errUnauthorized()
	}

	var out UserProfileOutput
	// すべて tx.Do の中に閉じ込める
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		up, err := uc.upRepo.FindByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}
		// レコードなし → 空 DTO でOK（HTTP層で 200 + null / {} / [] にするかは presenter で調整）
		if up == nil {
			out = UserProfileOutput{}
			return nil
		}

		// 年齢計算（例）
		var age *int
		if up.BirthDate != nil {
			a := util.CalcAge(up.BirthDate)
			age = &a
		}

		// Genderの辞書キーを組み立てる
		var genderLabelKey *string
		if up.Gender != nil && up.Gender.Code != "" {
			s := util.Fmt("master.gender.%s", up.Gender.Code)
			genderLabelKey = &s
		}

		// DTO 詰め替え
		out = UserProfileOutput{
			UserID:         up.UserID,
			FamilyName:     up.FamilyName,
			GivenName:      up.GivenName,
			FamilyNameKana: up.FamilyNameKana,
			GivenNameKana:  up.GivenNameKana,
			LegalName:      up.LegalName,
			LegalNameKana:  up.LegalNameKana,
			//BirthDate:      util.FormatDatePtr(up.BirthDate), // *time.Time -> *string にする util 想定
			BirthDate:      up.BirthDate, // *time.Time -> *string にする util 想定
			Age:            age,
			GenderID:       up.GenderID,
			GenderLabelKey: genderLabelKey,
			Initial:        up.Initial,
		}
		return nil
	}); err != nil {
		return UserProfileOutput{}, err
	}

	return out, nil
}
