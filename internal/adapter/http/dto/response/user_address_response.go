// Package response は HTTP レスポンス DTO を提供します。
// コントローラ → プレゼンター層で最終的に API が返す JSON 形を定義します。
package response

import "time"

// UserAddressResponse はユーザーの住所情報を表す構造体です。
// 各フィールドは住所エンティティに対応しており、API レスポンスやユースケースの出力 DTO として利用されます。
// null 許容の項目はポインタ型で表現され、未指定の場合は JSON に null が出力されます。
type UserAddressResponse struct {
	ID         uint64 `json:"id"` // 生成された主キー
	UserID     uint64 `json:"user_id"`
	PurposeID  uint64 `json:"purpose_id"`
	PurposeKey string `json:"purpose"`
	IsPrimary  bool   `json:"is_primary"` // 同一用途内での代表住所フラグ

	CountryCode        string  `json:"country_code"`        // ISO 3166-1 alpha-2 国コード
	AdministrativeArea *string `json:"administrative_area"` // 都道府県など上位行政区
	Locality           *string `json:"locality"`            // 市区町村
	DependentLocality  *string `json:"dependent_locality"`  // 町域・地区名など
	PostalCode         *string `json:"postal_code"`         // 郵便番号

	AddressLine1 string  `json:"address_line1"` // 番地・丁目
	AddressLine2 *string `json:"address_line2"` // 建物名・部屋番号
	AddressLine3 *string `json:"address_line3"` // 補足的な住所要素

	Latitude  *float64 `json:"latitude"`  // 緯度（任意）
	Longitude *float64 `json:"longitude"` // 経度（任意）

	CreatedAt time.Time `json:"createdAt"` // 登録日時（RFC3339）
	UpdatedAt time.Time `json:"updatedAt"` // 更新日時（RFC3339）
}
