// Package entity は住所の用途を表すドメインエンティティです
package entity

import "time"

// AddressPurpose は住所の用途を表すドメインエンティティです
type AddressPurpose struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Code        string    `gorm:"size:64;not null;unique;column:code"`
	DisplayName string    `gorm:"autoUpdateTime;column:display_name"`
	SortOrder   int       `gorm:"not null;default:0;column:sort_order"`
	IsActive    bool      `gorm:"not null;default:true;column:is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

// TableName は GORM のテーブル名を返します。
func (AddressPurpose) TableName() string {
	return "address_purposes"
}

// AddressPurposeID◯◯ は address_purposes テーブルの ID を表す定数です。
// マジックナンバーを避けるために、用途 ID を参照する箇所ではこれらを利用します。
const (
	// AddressPurposeIDHome は code=home の AddressPurpose のIDです。
	AddressPurposeIDHome uint64 = 1

	// AddressPurposeIDContact は code=contact の AddressPurpose のIDです。
	AddressPurposeIDContact uint64 = 2

	// AddressPurposeIDOffice は code=office の AddressPurpose のIDです。
	AddressPurposeIDOffice uint64 = 3

	// AddressPurposeIDShipping は code=shipping の AddressPurpose のIDです。
	AddressPurposeIDShipping uint64 = 4

	// AddressPurposeIDBilling は code=billing の AddressPurpose のIDです。
	AddressPurposeIDBilling uint64 = 5

	// AddressPurposeIDOther は code=other の AddressPurpose のIDです。
	AddressPurposeIDOther uint64 = 6
)

// AddressPurposeCode◯◯ は address_purposes テーブルの code を表す定数です。
// バリデーションや分岐で code を扱う場合には、文字列リテラルではなくこれらを利用します。
const (
	// AddressPurposeCodeHome は home 用途を表す code です。
	AddressPurposeCodeHome = "home"

	// AddressPurposeCodeContact は contact 用途を表す code です。
	AddressPurposeCodeContact = "contact"

	// AddressPurposeCodeOffice は office 用途を表す code です。
	AddressPurposeCodeOffice = "office"

	// AddressPurposeCodeShipping は shipping 用途を表す code です。
	AddressPurposeCodeShipping = "shipping"

	// AddressPurposeCodeBilling は billing 用途を表す code です。
	AddressPurposeCodeBilling = "billing"

	// AddressPurposeCodeOther は other 用途を表す code です。
	AddressPurposeCodeOther = "other"
)
