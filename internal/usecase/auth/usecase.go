// Package auth は、ユースケースのインターフェースを定義します。
package auth

import "context"

// Usecase は認証ユースケースの集合です。
type Usecase interface {

	// UpsertUserFromFirebase VerifyIDToken 後の upsert（users/auth_identities）と、ユーザー情報の返却。
	UpsertUserFromFirebase(ctx context.Context, in UpsertUserFromFirebaseInput) (AuthResultOutput, error)

	// 将来:
	// LinkProvider(ctx, in LinkProviderInput) (UserOutput, error)
	// UnlinkProvider(ctx, in UnlinkProviderInput) (UserOutput, error)
	// GetProfile(ctx context.Context, uid string) (UserOutput, error)
}
