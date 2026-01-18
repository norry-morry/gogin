// Package response は HTTP レスポンス DTO を提供します。
// コントローラ → プレゼンター層で最終的に API が返す JSON 形を定義します。
package response

import "time"

// UserEducationResponse はユーザーの学歴を表す構造体です。
// 各フィールドは学歴エンティティに対応しており、API レスポンスやユースケースの出力 DTO として利用されます。
// null 許容の項目はポインタ型で表現され、未指定の場合は JSON に null が出力されます。
type UserEducationResponse struct {
	ID                uint64    `json:"id"`
	UserID            uint64    `json:"user_id"`
	InstitutionName   string    `json:"institution_name"`
	FacultyName       *string   `json:"faculty_name"`
	DepartmentName    *string   `json:"department_name"`
	DegreeTypeID      uint64    `json:"degree_type_id"`
	DegreeType        string    `json:"degree_type"`
	EducationStatusID uint64    `json:"education_status_id"`
	EducationStatus   string    `json:"education_status"`
	EventDate         string    `json:"event_date"`
	Description       *string   `json:"description"`
	SortOrder         int       `json:"sort_order"`
	IsPublic          *bool     `json:"is_public"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	//EventDate         time.Time `json:"event_date"`
}
