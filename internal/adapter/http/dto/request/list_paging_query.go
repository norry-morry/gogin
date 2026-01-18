// Package request は HTTP ハンドラが受け取る入力DTO（Query/Body）の型をまとめます。
package request

// ListPagingQuery は 一覧系エンドポイントで共通に使うページングパラメータ。
type ListPagingQuery struct {
	// 1以上の整数。未指定なら 1。
	Page int `form:"page" binding:"omitempty,min=1"`

	// -1（全件）または 1..100。未指定なら 10。
	PerPage int `form:"per_page" binding:"omitempty"`

	// 並び順。asc/desc のみ。
	Order string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// Normalize はゼロ値にデフォルト値を入れたり、上限等を丸めます。
func (q *ListPagingQuery) Normalize() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PerPage == 0 {
		q.PerPage = 10
	}
	// -1 は「全件」。それ以外は 1..100 に丸める。
	if q.PerPage != -1 {
		if q.PerPage < 1 {
			q.PerPage = 10
			//} else if q.PerPage > 100 {
			//	q.PerPage = 100
		}
	}
	if q.Order == "" {
		q.Order = "desc"
	}
}
