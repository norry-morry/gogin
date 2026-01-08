package response

// UserProfileResponse は、ユーザーのプロフィール情報を API クライアントへ返却するためのレスポンス DTO です。
// ドメイン層の entity.UserProfile とは独立しており、API 仕様に基づいたフィールド構成と命名を持ちます。
// birth_date は ISO8601 形式へ整形され、age / age_group はユースケース層またはプレゼンタ層で算出されます。
// gender_label_key は多言語ラベル解決用のキーを表します。
type UserProfileResponse struct {
	UserID         uint64  `json:"user_id"`
	FamilyName     string  `json:"family_name"`
	GivenName      string  `json:"given_name"`
	FamilyNameKana string  `json:"family_name_kana"`
	GivenNameKana  string  `json:"given_name_kana"`
	LegalName      string  `json:"legal_name"`
	LegalNameKana  string  `json:"legal_name_kana,omitempty"`
	BirthDate      *string `json:"birth_date,omitempty"` // ISO8601形式などに整形
	Age            *int    `json:"age,omitempty"`        // ← usecase層で算出
	AgeGroup       *int    `json:"age_group,omitempty"`  // ← presenter層で算出
	GenderID       *uint8  `json:"gender_id"`
	GenderLabelKey *string `json:"gender_label_key"`
	Initial        *string `json:"initial,omitempty"`
}
