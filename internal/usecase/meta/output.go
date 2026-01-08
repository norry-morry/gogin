// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

// AddressPurposeDTO は、UI に返す単一の住所用途選択肢です。
type AddressPurposeDTO struct {
	ID        uint64 `json:"id"`
	Code      string `json:"code"`
	Label     string `json:"label"`
	Value     uint64 `json:"value"`
	Disabled  bool   `json:"disabled,omitempty"`
	SortOrder int    `json:"-"`
}

// ListAddressPurposeOutput は、住所用途の選択肢一覧の取得結果です。
type ListAddressPurposeOutput struct {
	Items []AddressPurposeDTO `json:"options"`
	ETag  string              `json:"-"`
}

// GenderDTO は、UI に返す単一の性別選択肢です。
type GenderDTO struct {
	ID        uint8  `json:"id"`
	Code      string `json:"code"`
	Label     string `json:"label"`
	Value     uint8  `json:"value"`
	IsActive  bool   `json:"is_active"`
	SortOrder uint8  `json:"sort_order"`
}

// ListGenderOutput は、UI に返す性別選択肢一覧の取得結果です。
type ListGenderOutput struct {
	Items []GenderDTO `json:"options"`
	ETag  string      `json:"-"`
}

// CountryDTO は、UI に返す単一の国選択肢です。
type CountryDTO struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ListCountryOutput は、UI に返す国選択肢一覧の取得結果です。
type ListCountryOutput struct {
	Items []CountryDTO `json:"countries"`
	ETag  string       `json:"-"`
}

// EducationStatusDTO は UIに返す単一の学歴状態の選択肢です
type EducationStatusDTO struct {
	ID        uint64 `json:"id"`
	Code      string `json:"code"`
	Label     string `json:"label"`
	Value     uint64 `json:"value"`
	Disabled  bool   `json:"disabled,omitempty"`
	SortOrder int    `json:"-"`
}

// ListEducationStatusOutput は 学歴状態の選択肢一覧の取得結果です
type ListEducationStatusOutput struct {
	Items []EducationStatusDTO `json:"options"`
	ETag  string               `json:"-"`
}

// DegreeTypeDTO は UIに返す学位種別の選択肢です
type DegreeTypeDTO struct {
	ID        uint64 `json:"id"`
	Code      string `json:"code"`
	Label     string `json:"label"`
	Value     uint64 `json:"value"`
	Disabled  bool   `json:"disabled,omitempty"`
	SortOrder int    `json:"-"`
}

// ListDegreeTypeOutput は 学位種別の選択肢一覧の取得結果です
type ListDegreeTypeOutput struct {
	Items []DegreeTypeDTO `json:"options"`
	ETag  string          `json:"-"`
}
