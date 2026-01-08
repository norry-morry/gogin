// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import "time"

// CreateInput は 学歴登録の入力DTOを表す
type CreateInput struct {
	UserID            uint64
	InstitutionName   string
	FacultyName       *string
	DepartmentName    *string
	DegreeTypeID      uint64
	EducationStatusID uint64
	EventDate         time.Time
	Description       *string
	IsPublic          *bool
}

// UpdateInput は 学歴更新の入力DTOを表す
type UpdateInput struct {
	EducationID       uint64
	UserID            uint64
	InstitutionName   string
	FacultyName       *string
	DepartmentName    *string
	DegreeTypeID      uint64
	EducationStatusID uint64
	EventDate         time.Time
	Description       *string
	IsPublic          *bool
}

// ListInput は 学歴一覧の入力DTOを表す
type ListInput struct {
	UserID uint64
}

// DeleteInput は 学歴削除の入力DTOを表す
type DeleteInput struct {
	EducationID uint64
	UserID      uint64
}

// ReorderInput は 学歴並び替えの入力DTOを表す
type ReorderInput struct {
	UserID       uint64
	EducationIDs []uint64
}

// HasEducationInput は 学歴有無の入力DTOを表す
type HasEducationInput struct {
	UserID uint64
}
