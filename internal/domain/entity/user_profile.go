// Package entity は ユーザーの履歴書・経歴書に記載する今人情報を表すドメインエンティティです。
package entity

import "time"

// UserProfile は ユーザーの履歴書・経歴書に記載する今人情報を表すドメインエンティティです。
type UserProfile struct {
	UserID         uint64 `json:"user_id" gorm:"primaryKey;column:user_id"`
	FamilyName     string `json:"family_name" gorm:"column:family_name;size:80;not null"`
	GivenName      string `json:"given_name"  gorm:"column:given_name;size:80;not null"`
	FamilyNameKana string `json:"family_name_kana" gorm:"column:family_name_kana;size:80;not null"`
	GivenNameKana  string `json:"given_name_kana"  gorm:"column:given_name_kana;size:80;not null"`

	// 生成列（STORED）。DBで作っているので読み取り専用にしておくと安全
	LegalName     string `json:"legal_name"       gorm:"column:legal_name;->"`      // CONCAT_WS(' ', family_name, given_name)
	LegalNameKana string `json:"legal_name_kana"  gorm:"column:legal_name_kana;->"` // CONCAT_WS(' ', family_name_kana, given_name_kana)

	BirthDate *time.Time `json:"birth_date,omitempty" gorm:"column:birth_date"`

	GenderID *uint8  `json:"gender_id,omitempty" gorm:"column:gender_id"`
	Gender   *Gender `gorm:"foreignKey:GenderID;references:ID"`
	Initial  *string `json:"initial,omitempty"   gorm:"column:initial;size:32"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName golang-migrate で auth_identities を作っている前提に合わせる
func (UserProfile) TableName() string {
	return "user_profiles"
}
