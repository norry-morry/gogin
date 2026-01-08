// Package entity は 学歴状態を表すドメインエンティティです
package entity

import "time"

// EducationStatus は 学歴状態を表すドメインエンティティです
type EducationStatus struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Code      string    `gorm:"size:32;not null;unique;column:code"`
	SortOrder int       `gorm:"not null;default:0;column:sort_order"`
	IsActive  bool      `gorm:"not null;default:true;column:is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

// TableName は GORMのテーブル名を返します
func (EducationStatus) TableName() string { return "education_statuses" }

// EducationStatus◯◯ は education_statuses テーブルの ID を表す定数です。
// マジックナンバーを避けるために、状態ID を参照する箇所ではこれらを利用します。
const (
	// EducationStatusEntrance は code=entrance の EducationStatus の IDです。-入学
	EducationStatusEntrance uint64 = 1

	// EducationStatusEnrolled は code=enrolled の EducationStatus のIDです。-在学中
	EducationStatusEnrolled uint64 = 2

	// EducationStatusLeaveOfAbsence は code=leave_of_absence の EducationStatus のIDです。-休学
	EducationStatusLeaveOfAbsence uint64 = 3

	// EducationStatusGraduated は code=graduated の EducationStatus のIDです。-卒業
	EducationStatusGraduated uint64 = 4

	// EducationStatusCompleted は code=completed の EducationStatus のIDです。-終了
	EducationStatusCompleted uint64 = 5

	// EducationStatusGraduationProspect は code=graduated の EducationStatus のIDです。-卒業見込み
	EducationStatusGraduationProspect uint64 = 6

	// EducationStatusWithdrawn は code=withdrawn の EducationStatus のIDです。-退学
	EducationStatusWithdrawn uint64 = 7

	// EducationStatusExpelled は code=withdrawn の EducationStatus のIDです。-除籍・放校
	EducationStatusExpelled uint64 = 8
)

// EducationStatus◯◯ は education_statuses テーブルの code を表す定数です。
// バリデーションや分岐で code を扱う場合には、文字列リテラルではなくこれらを利用します。
const (
	// EducationStatusCodeEntrance は 入学 を表すcodeです
	EducationStatusCodeEntrance = "entrance"

	// EducationStatusCodeEnrolled は 在学中 を表すcodeです
	EducationStatusCodeEnrolled = "enrolled"

	// EducationStatusCodeLeaveOfAbsence は 休学 を表すcodeです
	EducationStatusCodeLeaveOfAbsence = "leave_of_absence"

	// EducationStatusCodeGraduated は 卒業 を表すcodeです
	EducationStatusCodeGraduated = "graduated"

	// EducationStatusCodeCompleted は 修了 を表すcodeです
	EducationStatusCodeCompleted = "completed"

	// EducationStatusCodeGraduationProspect は 卒業見込 を表すcodeです
	EducationStatusCodeGraduationProspect = "graduation_prospect"

	// EducationStatusCodeWithdrawn は 退学 を表すcodeです
	EducationStatusCodeWithdrawn = "withdrawn"

	// EducationStatusCodeExpelled は 除籍・放校 を表すcodeです
	EducationStatusCodeExpelled = "expelled"
)
