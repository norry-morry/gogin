// Package valerr_test は、ドメイン層の値検証エラーを
// handler/presenter 向けの構造に変換するユーティリティ
// (internal/shared/valerr/domain.go) のユニットテストを提供する。
package valerr_test

import (
	"reflect"
	"testing"

	"resume/internal/domain/entity"
	"resume/internal/shared/valerr"
)

//
// ---------- FromDomain（最後のエラーだけ採用） ----------
//

// TestFromDomain_Empty は Problems が空の場合に nil が返ることを確認する。
func TestFromDomain_Empty(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{Problems: []entity.FieldError{}}
	got := valerr.FromDomain(inv)

	if got != nil {
		t.Errorf("expected nil, got=%v", got)
	}
}

// TestFromDomain_SingleProblem は 単一の FieldError が map に変換されることを確認する。
func TestFromDomain_SingleProblem(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{
		Problems: []entity.FieldError{
			{Field: "postalCode", Tag: "jppostal", Param: ""},
		},
	}

	got := valerr.FromDomain(inv)
	want := map[string]any{
		"postalCode": map[string]any{"tag": "jppostal", "param": ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestFromDomain_LastWins は 同一フィールドに複数エラーがある場合に最後の 1 件だけが採用されることを確認する。
func TestFromDomain_LastWins(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{
		Problems: []entity.FieldError{
			{Field: "administrativeArea", Tag: "required", Param: ""},
			{Field: "administrativeArea", Tag: "jpref", Param: ""},
		},
	}

	got := valerr.FromDomain(inv)
	want := map[string]any{
		"administrativeArea": map[string]any{"tag": "jpref", "param": ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

//
// ---------- FromDomainAll（すべてのエラーを配列で返す） ----------
//

// TestFromDomainAll_Empty は Problems が空の場合に nil が返ることを確認する。
func TestFromDomainAll_Empty(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{Problems: []entity.FieldError{}}
	got := valerr.FromDomainAll(inv)

	if got != nil {
		t.Errorf("expected nil, got=%v", got)
	}
}

// TestFromDomainAll_SingleProblem は 単一の FieldError が map に変換されることを確認する。
func TestFromDomainAll_SingleProblem(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{
		Problems: []entity.FieldError{
			{Field: "postalCode", Tag: "jppostal", Param: ""},
		},
	}

	got := valerr.FromDomainAll(inv)
	want := map[string]any{
		"postalCode": map[string]any{"tag": "jppostal", "param": ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestFromDomainAll_MultipleForSameField は 同一フィールドに複数エラーがある場合に配列として返されることを確認する。
func TestFromDomainAll_MultipleForSameField(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{
		Problems: []entity.FieldError{
			{Field: "postalCode", Tag: "len", Param: "7"},
			{Field: "postalCode", Tag: "jppostal", Param: ""},
		},
	}

	got := valerr.FromDomainAll(inv)
	want := map[string]any{
		"postalCode": []map[string]any{
			{"tag": "len", "param": "7"},
			{"tag": "jppostal", "param": ""},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v,\n got %v", want, got)
	}
}

// TestFromDomainAll_MixedThree は 1 件目は単体、2件目以降が追加されて配列になるケースを確認する。
func TestFromDomainAll_MixedThree(t *testing.T) {
	t.Parallel()

	inv := entity.InvalidAddressError{
		Problems: []entity.FieldError{
			{Field: "administrativeArea", Tag: "required", Param: ""},
			{Field: "administrativeArea", Tag: "jpref", Param: ""},
			{Field: "administrativeArea", Tag: "max_len", Param: "32"},
		},
	}

	got := valerr.FromDomainAll(inv)
	want := map[string]any{
		"administrativeArea": []map[string]any{
			{"tag": "required", "param": ""},
			{"tag": "jpref", "param": ""},
			{"tag": "max_len", "param": "32"},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v,\n got %v", want, got)
	}
}
