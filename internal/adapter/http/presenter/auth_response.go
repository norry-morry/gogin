// Package presenter defines response DTOs and presentation logic.
package presenter

import "resume/internal/domain/entity"

// AuthResponse は認証成功時のレスポンス DTO です
type AuthResponse struct {
	UID         string  `json:"uid"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhotoURL    *string `json:"photoUrl,omitempty"`
	Token       string  `json:"token,omitempty"` // アプリ側で独自セッショントークンを発行する場合
}

// FromDomainUserForAuth はログイン直後に返すユーザー情報を AuthResponse に変換します
func FromDomainUserForAuth(u *entity.User, token string) AuthResponse {
	return AuthResponse{
		UID:         u.UID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		PhotoURL:    u.PhotoURL,
		Token:       token,
	}
}
