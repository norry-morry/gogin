// Package auth は、認証関連ユースケースで使用する入力DTOを定義します。
package auth

// ProviderInfo は、Firebase の認証プロバイダ情報を表します。
// UpsertUserFromFirebaseInput 内の Providers フィールドで利用されます。
type ProviderInfo struct {
	Provider            string `validate:"required"`
	ProviderUserID      string `validate:"required"`
	ProviderDisplayName string `validate:"required"`
	EmailAtSignup       *string
}

// UpsertUserFromFirebaseInput は、VerifyIDToken 成功後に
// users / auth_identities テーブルへ登録・更新する際の入力DTOです。
// Firebase ユーザー情報と、紐付くプロバイダ情報を含みます。
type UpsertUserFromFirebaseInput struct {
	IDToken       string `validate:"required"`
	UID           string `validate:"required"`
	Email         *string
	EmailVerified bool
	DisplayName   *string
	PhotoURL      *string

	Providers []ProviderInfo
}
