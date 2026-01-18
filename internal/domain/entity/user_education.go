// Package entity はドメインエンティティを定義します。
package entity

import (
	"time"
)

// UserEducation はユーザーの学歴を表すドメインエンティティです
type UserEducation struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement;comment:学歴ID"`
	UserID            uint64    `gorm:"not null;index:idx_user_educations_user;index:idx_user_educations_user_sort,priority:1;comment:ユーザーID（users.id）"`
	InstitutionName   string    `gorm:"type:varchar(255);not null;comment:学校名"`
	FacultyName       *string   `gorm:"type:varchar(255);comment:学部"`
	DepartmentName    *string   `gorm:"type:varchar(255);comment:学科・専攻"`
	DegreeTypeID      uint64    `gorm:"comment:学位種別ID（degree_types.id）"`
	EducationStatusID uint64    `gorm:"not null;comment:学歴状態ID（education_statuses.id）"`
	EventDate         time.Time `gorm:"type:date;not null;comment:年月"`
	Description       *string   `gorm:"type:text;comment:補足"`
	SortOrder         int       `gorm:"not null;default:0;index:idx_user_educations_user_sort,priority:2;comment:表示順"`
	IsPublic          *bool     `gorm:"not null;default:true;comment:公開可否"`
	CreatedAt         time.Time `gorm:"autoCreateTime;comment:作成日時"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime;comment:更新日時"`

	// 関連
	DegreeType      *DegreeType      `gorm:"foreignKey:DegreeTypeID"`
	EducationStatus *EducationStatus `gorm:"foreignKey:EducationStatusID"`
}

// TableName は GORM のテーブル名を返します。
func (UserEducation) TableName() string {
	return "user_educations"
}
