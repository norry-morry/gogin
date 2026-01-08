package user

import "time"

// UserOutput はユーザー取得の出力DTOです。
// revive:disable-next-line:exported
type UserOutput struct {
	ID            uint64
	UID           string
	Email         *string
	EmailVerified bool
	DisplayName   *string
	PhotoURL      *string
	Disabled      bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
