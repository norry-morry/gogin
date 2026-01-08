// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import (
	"resume/internal/adapter/http/dto/response"
	"resume/internal/domain/entity"
)

// CreateOutput は 学歴登録の出力DTOを表す
// entity.UserEducation を内包し、アプリケーション層から Presenter 層へ渡す
type CreateOutput struct {
	Education *entity.UserEducation
}

// UpdateOutput は 学歴更新の出力DTOを表す
// entity.UserEducation を内包し、アプリケーション層から Presenter 層へ渡す
type UpdateOutput struct {
	Education *entity.UserEducation
}

// ListOutput は 学歴一覧の出力DTOを表す
// entity.UserEducationResponse を内包し、アプリケーション層から Presenter 層へ渡す
type ListOutput struct {
	Items []response.UserEducationResponse `json:"items"`
}

// DeleteOutput は 学歴削除の出漁DTOを表す
type DeleteOutput struct{}

// ReorderOutput は 学歴並び替えの出力DTOを表す
// 並び替え処理は成否を返すだけなので、中身は空
type ReorderOutput struct{}

// HasEducationOutput は 学歴有無の出力DTOを表す
type HasEducationOutput struct {
	Exists bool `json:"exists"`
}
