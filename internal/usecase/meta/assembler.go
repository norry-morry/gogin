// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"

	"resume/internal/domain/entity"
)

// assembleAddressPurpose は、住所用途エンティティの配列と「コード→表示名」を解決する関数を受け取り、
// DTO のスライスと ETag を組み立てます。resolve は (code, fallbackDisplayName) -> label を返す関数です。
func assembleAddressPurpose(
	list []entity.AddressPurpose,
	resolve func(code, fallback string) string,
) ListAddressPurposeOutput {

	opts := make([]AddressPurposeDTO, 0, len(list))
	for _, ap := range list {
		label := ap.Code
		if resolve != nil {
			label = resolve(ap.Code, ap.DisplayName)
		}
		opts = append(opts, AddressPurposeDTO{
			ID:        ap.ID,
			Code:      ap.Code,
			Label:     label,
			Value:     ap.ID,
			Disabled:  false,
			SortOrder: ap.SortOrder,
		})
	}

	// 並び順: sort_order ASC, value ASC
	sort.SliceStable(opts, func(i, j int) bool {
		if opts[i].SortOrder == opts[j].SortOrder {
			return opts[i].Code < opts[j].Code
		}
		return opts[i].SortOrder < opts[j].SortOrder
	})

	// 弱い ETag（value のみでOK）
	h := sha1.New()
	for _, o := range opts {
		_, _ = h.Write([]byte(o.Code))
		_, _ = h.Write([]byte{0})
	}
	etag := `W/"` + hex.EncodeToString(h.Sum(nil)) + `"`

	return ListAddressPurposeOutput{
		Items: opts,
		ETag:  etag,
	}
}

// assembleGender は、性別エンティティの配列と「コード→表示名」を解決する関数を受け取り、
// DTO のスライスと ETag を組み立てます。resolve は (code, fallbackDisplayName) -> label を返す関数です。
func assembleGender(
	list []entity.Gender,
	resolve func(code, fallback string) string,
) ListGenderOutput {
	opts := make([]GenderDTO, 0, len(list))
	for _, gn := range list {
		label := gn.Code
		if resolve != nil {
			label = resolve(gn.Code, gn.Code)
		}
		opts = append(opts, GenderDTO{
			ID:        gn.ID,
			Code:      gn.Code,
			Label:     label,
			Value:     gn.ID,
			IsActive:  gn.IsActive,
			SortOrder: gn.SortOrder,
		})
	}

	// 並び順: sort_order ASC, value ASC
	sort.SliceStable(opts, func(i, j int) bool {
		if opts[i].SortOrder == opts[j].SortOrder {
			return opts[i].Code < opts[j].Code
		}
		return opts[i].SortOrder < opts[j].SortOrder
	})

	// 弱い ETag（value のみでOK）
	h := sha1.New()
	for _, o := range opts {
		_, _ = h.Write([]byte(o.Code))
		_, _ = h.Write([]byte{0})
	}
	etag := `W/"` + hex.EncodeToString(h.Sum(nil)) + `"`

	return ListGenderOutput{
		Items: opts,
		ETag:  etag,
	}
}

// assembleCountry は、国エンティティの配列と「コード→表示名」を解決する関数を受け取り、
// DTO のスライスと ETag を組み立てます。resolve は (code, fallbackDisplayName) -> label を返す関数です。
func assembleCountry(
	list []entity.Country,
	resolve func(code, fallback string) string,
) ListCountryOutput {
	opts := make([]CountryDTO, 0, len(list))
	for _, cn := range list {
		label := cn.Code
		if resolve != nil {
			label = resolve(cn.Code, cn.Code)
		}
		opts = append(opts, CountryDTO{
			Value: cn.Code,
			Label: label,
		})
	}

	h := sha1.New()
	for _, o := range opts {
		_, _ = h.Write([]byte(o.Value))
		_, _ = h.Write([]byte{0})
	}
	etag := `W/"` + hex.EncodeToString(h.Sum(nil)) + `"`

	return ListCountryOutput{
		Items: opts,
		ETag:  etag,
	}
}

// assembleEducationStatus は、学歴状態エンティティの配列と「コード→表示名」を解決する関数を受け取り、
// DTO のスライスと ETag を組み立てます。resolve は (code, fallbackDisplayName) -> label を返す関数です。
func assembleEducationStatus(
	list []entity.EducationStatus,
	resolve func(code, fallback string) string,
) ListEducationStatusOutput {
	opts := make([]EducationStatusDTO, 0, len(list))
	for _, es := range list {
		label := es.Code
		if resolve != nil {
			label = resolve(es.Code, es.Code)
		}
		opts = append(opts, EducationStatusDTO{
			ID:        es.ID,
			Code:      es.Code,
			Label:     label,
			Value:     es.ID,
			Disabled:  false,
			SortOrder: es.SortOrder,
		})
	}

	sort.SliceStable(opts, func(i, j int) bool {
		if opts[i].SortOrder == opts[j].SortOrder {
			return opts[i].Code < opts[j].Code
		}
		return opts[i].SortOrder < opts[j].SortOrder
	})

	h := sha1.New()
	for _, o := range opts {
		_, _ = h.Write([]byte(o.Code))
		_, _ = h.Write([]byte{0})
	}
	etag := `W/"` + hex.EncodeToString(h.Sum(nil)) + `"`
	return ListEducationStatusOutput{
		Items: opts,
		ETag:  etag,
	}
}

// assembleDegreeType は、学位種別エンティティの配列と「コード→表示名」を解決する関数を受け取り、
// DTO のスライスと ETag を組み立てます。resolve は (code, fallbackDisplayName) -> label を返す関数です。
func assembleDegreeType(
	list []entity.DegreeType,
	resolve func(code, fallback string) string,
) ListDegreeTypeOutput {
	opts := make([]DegreeTypeDTO, 0, len(list))
	for _, dt := range list {
		label := dt.Code
		if resolve != nil {
			label = resolve(dt.Code, dt.Code)
		}
		opts = append(opts, DegreeTypeDTO{
			ID:        dt.ID,
			Code:      dt.Code,
			Label:     label,
			Value:     dt.ID,
			Disabled:  false,
			SortOrder: dt.SortOrder,
		})
	}
	sort.SliceStable(opts, func(i, j int) bool {
		if opts[i].SortOrder == opts[j].SortOrder {
			return opts[i].Code < opts[j].Code
		}
		return opts[i].SortOrder < opts[j].SortOrder
	})

	h := sha1.New()
	for _, o := range opts {
		_, _ = h.Write([]byte(o.Code))
		_, _ = h.Write([]byte{0})
	}
	etag := `W/"` + hex.EncodeToString(h.Sum(nil)) + `"`
	return ListDegreeTypeOutput{
		Items: opts,
		ETag:  etag,
	}
}
