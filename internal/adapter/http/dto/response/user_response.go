// Package response は HTTP レスポンス DTO を提供します。
// コントローラ → プレゼンター層で最終的に API が返す JSON 形を定義します。
package response

// User はユーザー情報をクライアントへ返すためのレスポンスDTOです。
// ドメインの entity.User とは分離されており、API 仕様に合わせたフィールド構成/命名を持ちます。
type User struct {
	ID            uint64  `json:"id"`
	UID           string  `json:"UID"`
	Email         *string `json:"email"`
	EmailVerified bool    `json:"emailVerified"`
	DisplayName   *string `json:"displayName"`
	PhotoURL      *string `json:"photoURL"`
	Disabled      bool    `json:"disabled"`
	LastLoginAt   *string `json:"lastLoginAt"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
	DeletedAt     *string `json:"deletedAt"`
}

// FromUsecase はユーザー情報をクライアントへ返すためのレスポンスDTOです。
func FromUsecase(uID uint64, name, email *string) User {
	return User{ID: uID, DisplayName: name, Email: email}
}
