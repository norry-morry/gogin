// Package rulestruct_test は、構造体レベルのカスタムバリデーション
// CoordPair のユニットテストを提供します。
package rulestruct_test

import (
	"testing"

	"github.com/go-playground/validator/v10"

	rulestruct "resume/internal/adapter/validation/rules/struct"
)

// locationPayload は CoordPair 検証の対象となるテスト用構造体です。
// ※ struct-level バリデーションなのでフィールドタグは不要。
type locationPayload struct {
	Latitude  float64
	Longitude float64
}

// newValidatorWithCoordPair は CoordPair を struct-level バリデーションとして
// 登録した validator.Validate を返します。
func newValidatorWithCoordPair() *validator.Validate {
	v := validator.New()
	v.RegisterStructValidation(rulestruct.CoordPair, locationPayload{})
	return v
}

// TestCoordPair_BothZero_NoError は、Latitude/Longitude が両方ゼロ値のときに
// エラーにならないことを確認します（＝両方未指定扱い）。
func TestCoordPair_BothZero_NoError(t *testing.T) {
	t.Parallel()

	v := newValidatorWithCoordPair()

	in := locationPayload{
		Latitude:  0,
		Longitude: 0,
	}

	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error when both latitude and longitude are zero, got=%v", err)
	}
}

// TestCoordPair_BothNonZero_NoError は、Latitude/Longitude が両方指定されているときに
// エラーにならないことを確認します。
func TestCoordPair_BothNonZero_NoError(t *testing.T) {
	t.Parallel()

	v := newValidatorWithCoordPair()

	in := locationPayload{
		Latitude:  35.0,
		Longitude: 135.0,
	}

	if err := v.Struct(in); err != nil {
		t.Fatalf("expected no error when both latitude and longitude are non-zero, got=%v", err)
	}
}

// TestCoordPair_OnlyLatitude_Set_ErrorOnBoth は、Latitude のみ指定されているときに
// 両フィールドに coordpair エラーが付くことを確認します。
func TestCoordPair_OnlyLatitude_Set_ErrorOnBoth(t *testing.T) {
	t.Parallel()

	v := newValidatorWithCoordPair()

	in := locationPayload{
		Latitude:  35.0,
		Longitude: 0,
	}

	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error when only latitude is set, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}
	if len(verrs) != 2 {
		t.Fatalf("expected 2 errors (latitude and longitude), got %d", len(verrs))
	}

	for _, e := range verrs {
		if e.Tag() != "coordpair" {
			t.Errorf("expected Tag=coordpair, got=%s", e.Tag())
		}
		if e.Param() != "both_or_none" {
			t.Errorf("expected Param=both_or_none, got=%s", e.Param())
		}
	}
}

// TestCoordPair_OnlyLongitude_Set_ErrorOnBoth は、Longitude のみ指定されているときに
// 両フィールドに coordpair エラーが付くことを確認します。
func TestCoordPair_OnlyLongitude_Set_ErrorOnBoth(t *testing.T) {
	t.Parallel()

	v := newValidatorWithCoordPair()

	in := locationPayload{
		Latitude:  0,
		Longitude: 135.0,
	}

	err := v.Struct(in)
	if err == nil {
		t.Fatalf("expected error when only longitude is set, got nil")
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}
	if len(verrs) != 2 {
		t.Fatalf("expected 2 errors (latitude and longitude), got %d", len(verrs))
	}

	for _, e := range verrs {
		if e.Tag() != "coordpair" {
			t.Errorf("expected Tag=coordpair, got=%s", e.Tag())
		}
		if e.Param() != "both_or_none" {
			t.Errorf("expected Param=both_or_none, got=%s", e.Param())
		}
	}
}
