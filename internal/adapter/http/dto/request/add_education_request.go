// Package request は HTTP ハンドラが受け取る入力DTO（Query/Body）の型をまとめます。
package request

// AddEducationRequest は 学歴登録ハンドラの作成リクエスト
type AddEducationRequest struct {
	InstitutionName   string  `json:"institution_name" binding:"required"`            // 学校名
	FacultyName       *string `json:"faculty_name" binding:"omitempty"`               // 学部
	DepartmentName    *string `json:"department_name" binding:"omitempty"`            // 学科・専攻
	DegreeTypeID      uint64  `json:"degree_type_id" binding:"required"`              // 学位種別ID
	EducationStatusID uint64  `json:"education_status_id" binding:"required"`         // 学歴状態ID
	EventDate         *string `json:"event_date" binding:"required,datetime=2006-01"` // 年月
	Description       *string `json:"description" binding:"omitempty"`                // 補足
	IsPublic          *bool   `json:"is_public" binding:"required"`                   // 公開・非公開
}

// Normalize は ハンドラ層でbind後に使う
func (r *AddEducationRequest) Normalize() {}
