// Package entity は 性別のドメインエンティティを定義します
package entity

import "time"

// Gender は性別マスタのエンティティです。
// コード値（male/female/other/unspecified）を主体に、
// アクティブ状態や並び順などを持ちます。
// 表示ラベルはi18n辞書で解決されるため、この構造体には含めません。
type Gender struct {
	ID        uint8     `json:"id" gorm:"primaryKey;column:id"` // TINYINT UNSIGNED に対応
	Code      string    `json:"code" gorm:"column:code;size:32;not null;unique"`
	SortOrder uint8     `json:"sort_order" gorm:"column:sort_order;not null;default:0"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName は GORM 用の明示的なテーブル名指定。
// （GORM のデフォルト規則では複数形にされるが、明示しておくと安全）
func (Gender) TableName() string {
	return "genders"
}

// GenderMaleID◯◯ は genders テーブルの ID を表す定数です。
// マジックナンバーを避けるために、性別 ID を参照する箇所ではこれらを利用します。
const (
	// GenderMaleID は code=home の Gender のIDです。
	GenderMaleID uint8 = 1

	// GenderFemaleID は code=home の Gender のIDです。
	GenderFemaleID uint8 = 2

	// GenderOtherID は code=home の Gender のIDです。
	GenderOtherID uint8 = 3
)

// GenderMaleCode◯◯ は genders テーブルの code を表す定数です。
// バリデーションや分岐で code を扱う場合には、文字列リテラルではなくこれらを利用します。
const (
	// GenderMaleCode は code=home 用途を表す code です。
	GenderMaleCode = "male"

	// GenderFemaleCode は code=female 用途を表す code です。
	GenderFemaleCode = "female"

	// GenderOtherCode は code=other 用途を表す code です。
	GenderOtherCode = "other"
)
