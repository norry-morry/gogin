// Package presenter は
package presenter

import "resume/internal/domain/entity"

// UserResponse は API レスポンス用の DTO（camelCase）です
type UserResponse struct {
	ID          uint64  `json:"id"`
	UID         string  `json:"uid"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhotoURL    *string `json:"photoURL,omitempty"`
	Disabled    bool    `json:"disabled"`
}

// FromDomainUser は domain.User を UserResponse に変換します
func FromDomainUser(u *entity.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		UID:         u.UID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		PhotoURL:    u.PhotoURL,
		Disabled:    u.Disabled,
	}
}

// FromDomainUsers は複数の domain.User を []UserResponse に変換します
func FromDomainUsers(users []*entity.User) []UserResponse {
	res := make([]UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, FromDomainUser(u))
	}
	return res
}
