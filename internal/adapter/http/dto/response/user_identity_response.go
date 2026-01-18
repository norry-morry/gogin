package response

// UserIdentity は、ユーザーに紐づく外部認証アカウントを API レスポンスとして表現する DTO です。
// フィールド名は Go の慣例に合わせて ID を使用し、JSON は snake_case で返します。
type UserIdentity struct {
	ID                  uint64  `json:"id"`
	UserID              uint64  `json:"user_id"`
	Provider            string  `json:"provider"`
	ProviderUserID      string  `json:"provider_user_id"`
	ProviderDisplayName string  `json:"provider_display_name"`
	EmailAtSignup       *string `json:"email_at_signup"`
	CreatedAt           string  `json:"created_at"`
}
