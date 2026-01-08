// Package profile は、ユーザープロフィールに関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package profile

import "context"

// HasUserProfile は、指定されたユーザーIDに紐づくプロフィール情報の有無を取得します。
// トランザクション内で UserProfileRepository を参照し、真偽値の DTO を返します。
func (uc *Interactor) HasUserProfile(ctx context.Context, in UserInput) (HasUserProfileOutput, error) {
	if in.UserID == 0 {
		return HasUserProfileOutput{}, errUnauthorized()
	}

	var exists bool
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {

		up, err := uc.upRepo.FindByUserID(txCtx, in.UserID)
		if err != nil {
			return err // ← エラーの場合は返して OK
		}

		exists = up != nil

		return nil // ← ココが重要
	}); err != nil {
		return HasUserProfileOutput{}, err
	}

	return HasUserProfileOutput{
		Exists: exists,
	}, nil
}
