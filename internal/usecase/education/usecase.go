// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import "context"

// Usecase は 学歴に関するユースケースの操作を定義します
type Usecase interface {
	// CreateUserEducation は 学歴を作成します
	CreateUserEducation(ctx context.Context, in CreateInput) (CreateOutput, error)

	// HasUserEducation は 指定ユーザーが学歴の登録をしているか判定します
	HasUserEducation(ctx context.Context, in HasEducationInput) (HasEducationOutput, error)

	// ListUserEducation は 指定ユーザーの学歴リストを取得します
	ListUserEducation(ctx context.Context, in ListInput) (ListOutput, error)

	// UpdateUserEducation は 学歴を更新します
	UpdateUserEducation(ctx context.Context, in UpdateInput) (UpdateOutput, error)

	// DeleteUserEducation は 学歴を物削します
	DeleteUserEducation(ctx context.Context, in DeleteInput) (DeleteOutput, error)

	// ReorderEducation は 学歴を並び替えします
	ReorderEducation(ctx context.Context, in ReorderInput) (ReorderOutput, error)
}
