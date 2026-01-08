package request

// CreateAddressRequest は住所の作成リクエスト。
type CreateAddressRequest struct {
	// 目的ID
	PurposeID *uint64 `json:"purpose_id" binding:"required,gt=0"`

	// ISO Alpha-2 を想定。厳密チェックはしないなら len=2 だけにして Normalize(大文字化)で吸収。
	CountryCode string `json:"country_code" binding:"required,len=2"`

	// JP のとき都道府県を必須にしたい → WhenJP(構造体検証)で担保するのでここは任意 + jpref
	AdministrativeArea *string `json:"administrative_area" binding:"omitempty,jp_pref,max=128"`

	// 任意。必要なら max のみ
	Locality          *string `json:"locality" binding:"omitempty,max=128"`
	DependentLocality *string `json:"dependent_locality" binding:"omitempty,max=128"`

	// JP では任意 or 必須かは要件次第。ここでは任意 + 書式は jp_postal に委ねる
	PostalCode *string `json:"postal_code" binding:"omitempty,jp_postal"`

	// 必須。前後スペース禁止 + 長さ
	AddressLine1 string  `json:"address_line1" binding:"required,max=160,no_edge_spaces"`
	AddressLine2 *string `json:"address_line2" binding:"omitempty,max=160,no_edge_spaces"`
	AddressLine3 *string `json:"address_line3" binding:"omitempty,max=160,no_edge_spaces"`

	// lat/lng は **両方セット or 両方未指定** → CoordPair(構造体検証)で担保
	Latitude  *float64 `json:"latitude" binding:"omitempty"`
	Longitude *float64 `json:"longitude" binding:"omitempty"`
}

// Normalize はハンドラ層で Bind 後に呼ぶ
func (r *CreateAddressRequest) Normalize() {
	// jpstring を使って正規化（import は省略）
	// r.CountryCode = strings.ToUpper(strings.TrimSpace(r.CountryCode))
	// ...必要な項目に jpstring.NormalizeJP / NormalizeDigitsHyphen を適用
}
