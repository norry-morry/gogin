// Package request は HTTP ハンドラが受け取る入力DTO（Query/Body）の型をまとめます。
package request

// ListIdentitiesQuery は /profile/identities のクエリパラメータ。
type ListIdentitiesQuery struct {
	// ページング共通部
	ListPagingQuery

	// 並び替えキー。許容値は sort_key バリデータで制御。
	// 例: created_at / provider / uid / email_at_signup
	Sort string `form:"sort" binding:"omitempty,sort_key=identities"`

	// フリーテキスト検索。未指定なら空扱い。
	Q *string `form:"q" binding:"omitempty,max=200"`
}

// Normalize はゼロ値にデフォルト値を入れたり、上限等を丸めます。
func (q *ListIdentitiesQuery) Normalize() {
	// ページング共通処理
	q.ListPagingQuery.Normalize()

	// ソートキーのデフォルト
	if q.Sort == "" {
		q.Sort = "created_at"
	}
}
