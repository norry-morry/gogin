// Package entity はユーザー関連のドメインエンティティを定義します。
package entity

import "time"

// User はユーザーを表すドメインエンティティです。
type User struct {
	ID            uint64  `json:"id" gorm:"primary_key"`
	UID           string  `json:"uid"`
	Email         *string `json:"email" gorm:"uniqueIndex"`
	EmailVerified bool    `json:"email_verified"`
	DisplayName   *string `json:"display_name"`
	PhotoURL      *string `json:"photo_url"`
	Disabled      bool    `json:"disabled"`
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

// TableName golang-migrate で users を作っている前提に合わせる
func (User) TableName() string {
	return "users"
}
