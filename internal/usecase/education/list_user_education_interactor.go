// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import "context"

// ListUserEducation は 指定ユーザーの学歴リストを取得します
func (uc *Interactor) ListUserEducation(
	ctx context.Context,
	in ListInput,
) (ListOutput, error) {
	if in.UserID == 0 {
		return ListOutput{}, errUnauthorized()
	}
	educations, err := uc.ueRepo.ListByUserID(ctx, in.UserID)
	if err != nil {
		return ListOutput{}, err
	}
	return ListOutput{
		Items: ToUserEducationResponses(educations),
	}, nil
}
