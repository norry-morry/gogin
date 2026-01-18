// Package request は HTTP ハンドラが受け取る入力DTO（Query/Body）の型をまとめます。
package request

// ListUserAddressQuery は /profile/address のクエリパラメータ。
type ListUserAddressQuery struct {
	ListPagingQuery
	Sort      string  `form:"sort" binding:"omitempty,sort_key=addresses"`
	PurposeID *uint64 `form:"purpose_id" binding:"omitempty"`
	Country   *string `form:"country" binding:"omitempty"`
	City      *string `form:"city" binding:"omitempty"`
}

// Normalize はゼロ値にデフォルト値を入れたり、上限等を丸めます。
func (q *ListUserAddressQuery) Normalize() {
	q.ListPagingQuery.Normalize()

	if q.Sort == "" {
		q.Sort = "created_at"
	}
}
