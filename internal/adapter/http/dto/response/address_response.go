package response

import "time"

// Address はユーザーの住所情報を表す構造体です。
// 各フィールドは住所エンティティに対応しており、API レスポンスやユースケースの出力 DTO として利用されます。
// null 許容の項目はポインタ型で表現され、未指定の場合は JSON に null が出力されます。
type Address struct {
	ID        uint64 `json:"id"`        // 生成された主キー
	PurposeID uint64 `json:"purposeId"` // 用途 (例: 請求先, 居住地 など)
	IsPrimary bool   `json:"isPrimary"` // 同一用途内での代表住所フラグ

	CountryCode        string  `json:"countryCode"`        // ISO 3166-1 alpha-2 国コード
	AdministrativeArea *string `json:"administrativeArea"` // 都道府県など上位行政区
	Locality           *string `json:"locality"`           // 市区町村
	DependentLocality  *string `json:"dependentLocality"`  // 町域・地区名など
	PostalCode         *string `json:"postalCode"`         // 郵便番号

	AddressLine1 string  `json:"addressLine1"` // 番地・丁目
	AddressLine2 *string `json:"addressLine2"` // 建物名・部屋番号
	AddressLine3 *string `json:"addressLine3"` // 補足的な住所要素

	Latitude  *float64 `json:"latitude"`  // 緯度（任意）
	Longitude *float64 `json:"longitude"` // 経度（任意）

	CreatedAt time.Time `json:"createdAt"` // 登録日時（RFC3339）
	UpdatedAt time.Time `json:"updatedAt"` // 更新日時（RFC3339）
}
