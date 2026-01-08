// Package entity は 学位種別を表すドメインエンティティです
package entity

import "time"

// DegreeType は 学位状態を表すドメインエンティティです
type DegreeType struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Code      string    `gorm:"size:64;not null;unique;column:code"`
	SortOrder int       `gorm:"not null;default:0;column:sort_order"`
	IsActive  bool      `gorm:"not null;default:true;column:is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

// TableName は GORM のテーブル名を返します。
func (DegreeType) TableName() string {
	return "degree_types"
}

// DegreeType◯◯ は degree_types テーブルの ID を表す定数です。
// マジックナンバーを避けるために、用途 ID を参照する箇所ではこれらを利用します。
const (
	// DegreeTypeHighSchool は code=high_school の DegreeType のIDです。
	DegreeTypeHighSchool uint64 = 1

	// DegreeTypeVocational は code=vocational の DegreeType のIDです。
	DegreeTypeVocational uint64 = 2

	// DegreeTypeJuniorCollege は code=junior_college の DegreeType のIDです。
	DegreeTypeJuniorCollege uint64 = 3

	// DegreeTypeBachelor は code=bachelor の DegreeType のIDです。
	DegreeTypeBachelor uint64 = 4

	// DegreeTypeMaster は code=master の DegreeType のIDです。
	DegreeTypeMaster uint64 = 5

	// DegreeTypeDoctor は code=doctor の DegreeType のIDです。
	DegreeTypeDoctor uint64 = 6

	// DegreeTypeOther は code=other の DegreeType のIDです。
	DegreeTypeOther uint64 = 7
)

// DegreeType◯◯ は degree_types テーブルの code を表す定数です。
// バリデーションや分岐で code を扱う場合には、文字列リテラルではなくこれらを利用します。
const (
	// DegreeTypeCodeHighSchool は code=high_school を表す code です
	DegreeTypeCodeHighSchool = "highSchool"

	// DegreeTypeCodeVocational は code=vocational を表す code です
	DegreeTypeCodeVocational = "vocational"

	// DegreeTypeCodeJuniorCollege は code=junior_college を表す code です
	DegreeTypeCodeJuniorCollege = "juniorCollege"

	// DegreeTypeCodeBachelor は code=bachelor を表す code です
	DegreeTypeCodeBachelor = "bachelor"

	// DegreeTypeCodeMaster は code=master を表す code です
	DegreeTypeCodeMaster = "master"

	// DegreeTypeCodeDoctor は code=doctor を表す code です
	DegreeTypeCodeDoctor = "doctor"

	// DegreeTypeCodeOther は code=other を表す code です
	DegreeTypeCodeOther = "other"
)
