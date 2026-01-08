// Package presenter は
package presenter

import (
	"time"

	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
)

// ToUserResponse はレスポンス DTO です
func ToUserResponse(e *entity.User) response.User {
	var lastLoginAtStr *string
	if e.LastLoginAt != nil {
		s := e.LastLoginAt.Format(time.RFC3339) // "2025-10-29T17:45:00Z" など
		lastLoginAtStr = &s
	}
	var deletedAtStr *string
	if e.DeletedAt != nil {
		d := e.DeletedAt.Format(time.RFC3339) // "2025-10-29T17:45:00Z" など
		deletedAtStr = &d
	}
	return response.User{
		ID:            e.ID,
		UID:           e.UID,
		Email:         e.Email,
		EmailVerified: e.EmailVerified,
		DisplayName:   e.DisplayName,
		PhotoURL:      e.PhotoURL,
		Disabled:      e.Disabled,
		LastLoginAt:   lastLoginAtStr,
		CreatedAt:     e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     e.UpdatedAt.Format(time.RFC3339),
		DeletedAt:     deletedAtStr,
	}
}
