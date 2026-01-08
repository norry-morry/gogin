// Package util_test は NoopMasker (internal/shared/util/noop_masker.go)
// のユニットテストを提供する。
package util_test

import (
	"reflect"
	"testing"

	"resume/internal/shared/util"
)

//
// ---------- MaskByKey ----------
//

// TestNoopMasker_MaskByKey_ReturnsSameValue は MaskByKey が入力値をそのまま返すことを確認する。
func TestNoopMasker_MaskByKey_ReturnsSameValue(t *testing.T) {
	t.Parallel()

	m := util.NoopMasker{}

	tests := []struct {
		name string
		in   any
	}{
		{name: "string", in: "secret"},
		{name: "int", in: 123},
		{name: "map", in: map[string]any{"key": "value"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := m.MaskByKey("any_key", tt.in)
			if !reflect.DeepEqual(got, tt.in) {
				t.Errorf("expected %#v, got %#v", tt.in, got)
			}
		})
	}
}

//
// ---------- MaskMapShallow ----------
//

// TestNoopMasker_MaskMapShallow_ReturnsSameMap は MaskMapShallow が引数の map をそのまま返すことを確認する。
func TestNoopMasker_MaskMapShallow_ReturnsSameMap(t *testing.T) {
	t.Parallel()

	m := util.NoopMasker{}
	src := map[string]any{
		"password": "secret",
		"email":    "user@example.com",
	}

	got := m.MaskMapShallow(src)

	// 内容が同じであること
	if !reflect.DeepEqual(got, src) {
		t.Fatalf("expected %#v, got %#v", src, got)
	}

	// map は参照型なので、noop 実装では同じ map インスタンスを返していることが期待される。
	// もし将来コピー実装になってもテストが壊れないように、ポインタ比較はログにとどめる。
	gotRef := m.MaskMapShallow(src)
	if &gotRef != &src {
		t.Log("MaskMapShallow may return a copy; for NoopMasker this is acceptable as long as contents are equal.")
	}
}

//
// ---------- MaskAnyRecursive ----------
//

// TestNoopMasker_MaskAnyRecursive_ReturnsSameValue は MaskAnyRecursive が入力値をそのまま返すことを確認する。
func TestNoopMasker_MaskAnyRecursive_ReturnsSameValue(t *testing.T) {
	t.Parallel()

	m := util.NoopMasker{}

	tests := []struct {
		name string
		in   any
	}{
		{name: "string", in: "secret"},
		{name: "slice", in: []string{"a", "b"}},
		{name: "nested map", in: map[string]any{"nested": map[string]any{"k": "v"}}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := m.MaskAnyRecursive("any_key", tt.in)
			if !reflect.DeepEqual(got, tt.in) {
				t.Errorf("expected %#v, got %#v", tt.in, got)
			}
		})
	}
}
