// Package rulestruct は、構造体レベルで適用されるカスタム検証ルールを提供します。
// 例えば、緯度と経度の両方が指定されているか／どちらも未指定であるかを検証します。
package rulestruct

import "github.com/go-playground/validator/v10"

// CoordPair は、緯度 (Latitude) と経度 (Longitude) の両方が
// 「両方指定」または「両方未指定」であることを検証します。
func CoordPair(sl validator.StructLevel) {
	lat := sl.Current().FieldByName("Latitude")
	lng := sl.Current().FieldByName("Longitude")
	if !lat.IsValid() || !lng.IsValid() {
		return
	}
	hasLat := !lat.IsZero()
	hasLng := !lng.IsZero()
	if hasLat != hasLng {
		sl.ReportError(lat.Interface(), "latitude", "latitude", "coordpair", "both_or_none")
		sl.ReportError(lng.Interface(), "longitude", "longitude", "coordpair", "both_or_none")
	}
}
