package entity

import "time"

// Country は国情報のドメインエンティティです。
type Country struct {
	Code        string    // ISO 3166-1 alpha-2 コード（例: JP, US）
	IsSupported bool      // システムで利用可能か
	SortOrder   int       // 並び順（小さい順に優先）
	PhoneCode   *string   // 国際電話コード（例: "+81"）
	CreatedAt   time.Time // 作成日時
	UpdatedAt   time.Time // 更新日時
}

// TableName は GORM 用の明示的なテーブル名指定。
// （GORM のデフォルト規則では複数形にされるが、明示しておくと安全）
func (Country) TableName() string {
	return "countries"
}

// IsActive は利用可能状態を返します。
func (c Country) IsActive() bool {
	return c.IsSupported
}

// DisplayCode は内部コード（大文字）を返します。
func (c Country) DisplayCode() string {
	return c.Code
}
