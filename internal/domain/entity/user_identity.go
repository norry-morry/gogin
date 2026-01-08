// Package entity はユーザーの外部プロバイダ連携を表すドメインエンティティです。
package entity

import "time"

// UserIdentity はユーザーの外部プロバイダ連携を表すドメインエンティティです。
type UserIdentity struct {
	ID                  uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID              uint64    `json:"user_id" gorm:"index:idx_auth_user"` // 逆引き用 index
	Provider            string    `json:"provider" gorm:"size:64;index:uq_provider_user,unique"`
	ProviderUserID      string    `json:"provider_user_id" gorm:"size:255;index:uq_provider_user,unique"`
	ProviderDisplayName string    `json:"provider_display_name" gorm:"size:255"`
	EmailAtSignup       *string   `json:"email_at_signup" gorm:"size:320"`
	CreatedAt           time.Time `json:"created_at"`
}

// TableName golang-migrate で auth_identities を作っている前提に合わせる
func (UserIdentity) TableName() string {
	return "user_identities"
}
