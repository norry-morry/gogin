package request

// OrderEducationRequest は 学歴の並び順変更のリクエスト
type OrderEducationRequest struct {
	EducationIDs []uint64 `json:"education_ids" binding:"required"`
}
