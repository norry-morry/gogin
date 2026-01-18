// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import "resume/internal/shared/apperr"

func errUnauthorized() error {
	return apperr.New(
		apperr.CodeUnprocessable,
		"unauthorized",
		nil,
	)
}
