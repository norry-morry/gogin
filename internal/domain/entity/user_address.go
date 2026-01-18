// Package entity はフリーランサーの住所を表すドメインエンティティです
package entity

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"resume/internal/shared/jpstring"
)

// UserAddress はフリーランサーの住所を表すドメインエンティティです
type UserAddress struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UserID    uint64 `gorm:"not null;index:idx_user_primary,priority:1"`               // users.id FK
	PurposeID uint64 `gorm:"not null;index:uq_user_purpose_primary,unique,priority:2"` // address_purposes.id FK

	IsPrimary bool `gorm:"not null;default:false;index:idx_user_primary,priority:2"`

	CountryCode        string  `gorm:"type:char(2);not null;index:idx_country_postal,priority:1"`
	AdministrativeArea *string `gorm:"size:128;index:idx_loc,priority:2"`
	Locality           *string `gorm:"size:128;index:idx_loc,priority:3"`
	DependentLocality  *string `gorm:"size:128"`
	PostalCode         *string `gorm:"size:32;index:idx_country_postal,priority:2"`
	SortingCode        *string `gorm:"size:32"`

	AddressLine1 string  `gorm:"size:160;not null"`
	AddressLine2 *string `gorm:"size:160"`
	AddressLine3 *string `gorm:"size:160"`

	Organization string `gorm:"size:160"`
	GivenName    string `gorm:"size:80"`
	FamilyName   string `gorm:"size:80"`

	LanguageCode string   `gorm:"size:35"`
	Latitude     *float64 `gorm:"type:decimal(9,6)"`
	Longitude    *float64 `gorm:"type:decimal(9,6)"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// 関連
	Purpose *AddressPurpose `gorm:"foreignKey:PurposeID;references:ID"`
}

// TableName は GORM のテーブル名を返します。
func (UserAddress) TableName() string {
	return "user_addresses"
}

// UserAddressParam は UserAddress を生成するための入力値を表します。
// バリデーションは NewUserAddress 内で行われ、無効な場合はエラーを返します。
type UserAddressParam struct {
	UserID, PurposeID                                                        uint64
	IsPrimary                                                                bool
	CountryCode                                                              string
	AdministrativeArea, Locality, DependentLocality, PostalCode, SortingCode *string
	AddressLine1                                                             string
	AddressLine2, AddressLine3                                               *string
	Latitude, Longitude                                                      *float64
}

// NewUserAddress は UserAddress エンティティを生成します。
// 与えられた UserAddressParam が無効な場合は ErrInvalidAddress などを返します。
func NewUserAddress(p UserAddressParam) (*UserAddress, error) {

	if (p.Latitude != nil && p.Longitude == nil) || (p.Latitude == nil && p.Longitude != nil) {
		return nil, errors.Join(ErrInvalidAddress, errCoordPair)
	}

	if p.Latitude != nil && (*p.Latitude < -90 || *p.Latitude > 90) {
		return nil, errors.Join(ErrInvalidAddress, errLatRange)
	}

	if p.Longitude != nil && (*p.Longitude < -180 || *p.Longitude > 180) {
		return nil, errors.Join(ErrInvalidAddress, errLngRange)
	}

	if p.CountryCode == "" {
		return nil, errors.Join(ErrInvalidAddress, errCountryCode)
	}

	if p.CountryCode == "JP" {
		if p.AdministrativeArea == nil || *p.AdministrativeArea == "" {
			return nil, errors.Join(ErrInvalidAddress, errJPPrefRequired)
		}
		if !jpstring.IsValidPrefecture(*p.AdministrativeArea) {
			return nil, errors.Join(ErrInvalidAddress, errJPPrefInvalid)
		}
	}

	if p.AddressLine1 == "" {
		return nil, errors.Join(ErrInvalidAddress, errAddressLine1Required)
	}

	var probs []FieldError

	// 標準化
	cc := strings.ToUpper(strings.TrimSpace(p.CountryCode))
	if len(cc) != 2 {
		probs = append(probs, FieldError{"countryCode", "len", "2"})
	}

	// addressLine1 必須（最終防衛）
	if strings.TrimSpace(p.AddressLine1) == "" {
		probs = append(probs, FieldError{"addressLine1", "required", ""})
	}

	// lat/lng のペア & range
	if (p.Latitude == nil) != (p.Longitude == nil) {
		probs = append(probs,
			FieldError{"latitude", "coordpair", "both_or_none"},
			FieldError{"longitude", "coordpair", "both_or_none"},
		)
	} else if p.Latitude != nil && p.Longitude != nil {
		if *p.Latitude < -90 || *p.Latitude > 90 {
			probs = append(probs, FieldError{"latitude", "range", "[-90,90]"})
		}
		if *p.Longitude < -180 || *p.Longitude > 180 {
			probs = append(probs, FieldError{"longitude", "range", "[-180,180]"})
		}
	}

	// JP 固有
	if cc == "JP" {
		// 都道府県必須 + 妥当性
		if p.AdministrativeArea == nil || strings.TrimSpace(*p.AdministrativeArea) == "" {
			probs = append(probs, FieldError{"administrativeArea", "jppref_required", "required_when:countryCode=JP"})
		} else {
			name := jpstring.NormalizeSpaces(*p.AdministrativeArea)
			if !jpstring.IsValidPrefecture(name) {
				probs = append(probs, FieldError{"administrativeArea", "jpref", ""})
			} else {
				// 正規化して反映（スペース整形）
				p.AdministrativeArea = &name
			}
		}
		// 郵便番号：存在すれば書式チェック（全角→半角・ハイフン統一後）
		if p.PostalCode != nil && strings.TrimSpace(*p.PostalCode) != "" {
			norm := jpstring.NormalizeDigitsHyphen(*p.PostalCode)
			if !rePostalJP.MatchString(norm) {
				probs = append(probs, FieldError{"postalCode", "jppostal", ""})
			} else {
				// 保存形を正規化に寄せたいならここで反映
				p.PostalCode = &norm
			}
		}
	}

	if len(probs) > 0 {
		return nil, InvalidAddressError{Problems: probs}
	}

	// 住所文字列の軽い正規化（任意。UI側ですでにやっていれば薄くてもOK）
	addr1 := jpstring.NormalizeJP(strings.TrimSpace(p.AddressLine1))
	var addr2, addr3 *string
	if p.AddressLine2 != nil {
		s := jpstring.NormalizeJP(strings.TrimSpace(*p.AddressLine2))
		addr2 = &s
	}
	if p.AddressLine3 != nil {
		s := jpstring.NormalizeJP(strings.TrimSpace(*p.AddressLine3))
		addr3 = &s
	}
	var locality, depLoc *string
	if p.Locality != nil {
		s := jpstring.NormalizeJP(strings.TrimSpace(*p.Locality))
		locality = &s
	}
	if p.DependentLocality != nil {
		s := jpstring.NormalizeJP(strings.TrimSpace(*p.DependentLocality))
		depLoc = &s
	}
	// ここで不変条件を最終チェック（例: len(CountryCode)==2、lat/lngのペア等）
	a := &UserAddress{
		UserID:    p.UserID,
		PurposeID: p.PurposeID,
		IsPrimary: p.IsPrimary,

		CountryCode:        cc, // ← cc を使う
		AdministrativeArea: p.AdministrativeArea,
		Locality:           locality, // ← 正規化版
		DependentLocality:  depLoc,   // ← 正規化版
		PostalCode:         p.PostalCode,
		SortingCode:        p.SortingCode,

		AddressLine1: addr1, // ← 正規化版
		AddressLine2: addr2, // ← 正規化版
		AddressLine3: addr3, // ← 正規化版

		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}
	return a, nil
}

var (
	rePostalJP = regexp.MustCompile(`^\d{3}-?\d{4}$`)

	// ErrInvalidAddress は住所エンティティの生成/検証に失敗したことを示すエラーです。
	// errors.Is(err, ErrInvalidAddress) で判定できます。
	ErrInvalidAddress = errors.New("invalid address")

	// ErrCoordPair は緯度と経度が「両方指定」または「両方未指定」の条件を満たしていないことを示すエラーです。
	// 片方のみ指定された場合に発生します。
	errCoordPair            = errors.New("latitude/longitude must be both set or both nil")
	errLatRange             = errors.New("latitude out of range")
	errLngRange             = errors.New("longitude out of range")
	errCountryCode          = errors.New("invalid country code")
	errJPPrefRequired       = errors.New("administrativeArea required for JP")
	errJPPrefInvalid        = errors.New("invalid Japanese prefecture")
	errAddressLine1Required = errors.New("addressLine1 required")
)

// FieldError は フィールド単位のエラー情報
type FieldError struct {
	Field string // JSON名で統一: "countryCode", "latitude" など
	Tag   string // "required", "coordpair", "range", "jppostal", "jpref_required" etc.
	Param string // 追加情報（"both_or_none", "[-90,90]", "required_when:countryCode=JP" 等）
}

// InvalidAddressError は住所が無効である場合の詳細を保持するエラー型です。
// フィールド単位の理由や原因となる値など、追加情報を含められます。
type InvalidAddressError struct {
	Problems []FieldError
}

func (e InvalidAddressError) Error() string { return ErrInvalidAddress.Error() }

// Is は errors.Is(e, target) の判定で ErrInvalidAddress と同一視できるようにします。
// これにより、呼び出し側は詳細型に依存せず「無効な住所」エラーを判定できます。
func (e InvalidAddressError) Is(target error) bool {
	return target == ErrInvalidAddress
}
