// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

// ListAddressPurposeInput は、住所用途の選択肢一覧を取得するユースケースの入力です。
type ListAddressPurposeInput struct {
	Locale          string
	IfNoneMatchEtag string
	IncludeInactive bool
	TenantID        uint64
}

// ListGenderInput は、性別の選択肢一覧を取得するユースケースの入力です。
type ListGenderInput struct {
	Locale string
	//IfNoneMatchEtag string
	//IncludeInactive bool
	//TenantID        uint64
}

// ListCountryInput は、国の選択肢一覧を取得するユースケースの入力です。
type ListCountryInput struct {
	Locale        string
	OnlySupported bool
	//IfNoneMatchEtag string
	//IncludeInactive bool
	//TenantID        uint64
}

// ListEducationStatusInput は 学歴状態一覧を取得するユースケースの入力です
type ListEducationStatusInput struct {
	Locale string
}

// ListDegreeTypeInput は 学位種別一覧を取得するユースケースの入力です
type ListDegreeTypeInput struct {
	Locale string
}
