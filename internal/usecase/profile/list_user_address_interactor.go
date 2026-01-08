// Package profile は、ユーザー住所に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package profile

import (
	"context"
	"log/slog"
)

// ListUserAddress は 認証済みユーザーの住所一覧を取得します。
func (uc *Interactor) ListUserAddress(ctx context.Context, in ListUserAddressInput) (ListUserAddressOutput, error) {
	if in.UserID == 0 {
		return ListUserAddressOutput{}, errUnauthorized()
	}
	slog.Debug("interactor", "ctx", ctx, "in", in)
	addresses, total, err := uc.uaRepo.ListByUserIDWithSpec(ctx, in.UserID, in)
	if err != nil {
		return ListUserAddressOutput{}, err
	}

	return ListUserAddressOutput{
		Items: ToUserAddressResponses(addresses),
		Total: total,
	}, nil
}
