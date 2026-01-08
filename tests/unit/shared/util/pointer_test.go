// Package util_test は pointer ユーティリティ
// (internal/shared/util/pointer.go) のユニットテストを提供する。
package util_test

import (
	"reflect"
	"testing"

	"resume/internal/shared/util"
)

//
// ---------- Clone ----------
//

// TestClone_Nil は nil ポインタを渡した場合に nil が返ることを確認する。
func TestClone_Nil(t *testing.T) {
	t.Parallel()

	var p *int
	got := util.Clone(p)
	if got != nil {
		t.Errorf("expected nil, got=%v", got)
	}
}

// TestClone_Value は値を持つポインタをクローンした場合に
// 値は同じだが別ポインタが返ることを確認する。
func TestClone_Value(t *testing.T) {
	t.Parallel()

	orig := 10
	p := &orig

	got := util.Clone(p)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}

	if *got != *p {
		t.Errorf("expected cloned value=%d, got=%d", *p, *got)
	}

	// アドレスが異なること（ディープコピー）
	if got == p {
		t.Errorf("expected different pointer (deep copy), but addresses are equal")
	}
}

// TestClone_Struct は構造体でも値がコピーされることを確認する。
func TestClone_Struct(t *testing.T) {
	t.Parallel()

	type user struct {
		Name string
		Age  int
	}

	orig := &user{Name: "Alice", Age: 20}
	got := util.Clone(orig)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}

	if !reflect.DeepEqual(orig, got) {
		t.Errorf("expected %#v, got %#v", orig, got)
	}

	if got == orig {
		t.Errorf("expected different pointer, but got same address")
	}
}

//
// ---------- ToPtr ----------
//

// TestToPtr_Primitive はプリミティブ型をポインタ化できることを確認する。
func TestToPtr_Primitive(t *testing.T) {
	t.Parallel()

	s := "hello"
	i := 42

	ps := util.ToPtr(s)
	if ps == nil || *ps != s {
		t.Errorf("expected %q, got=%v", s, ps)
	}

	pi := util.ToPtr(i)
	if pi == nil || *pi != i {
		t.Errorf("expected %d, got=%v", i, pi)
	}
}

// TestToPtr_Struct は構造体をポインタ化できることを確認する。
func TestToPtr_Struct(t *testing.T) {
	t.Parallel()

	type user struct {
		ID int
	}
	u := user{ID: 1}

	pu := util.ToPtr(u)
	if pu == nil {
		t.Fatal("expected non-nil pointer")
	}
	if pu.ID != u.ID {
		t.Errorf("expected %#v, got %#v", u, *pu)
	}
}

//
// ---------- Deref ----------
//

// TestDeref_NonNil は非 nil ポインタから値を取り出せることを確認する。
func TestDeref_NonNil(t *testing.T) {
	t.Parallel()

	i := 5
	got := util.Deref(&i)
	if got != i {
		t.Errorf("expected %d, got %d", i, got)
	}

	s := "world"
	gotS := util.Deref(&s)
	if gotS != s {
		t.Errorf("expected %q, got %q", s, gotS)
	}
}

// TestDeref_Nil は nil ポインタからゼロ値が返ることを確認する。
func TestDeref_Nil(t *testing.T) {
	t.Parallel()

	var pi *int
	gotI := util.Deref(pi)
	if gotI != 0 {
		t.Errorf("expected 0, got %d", gotI)
	}

	var ps *string
	gotS := util.Deref(ps)
	if gotS != "" {
		t.Errorf("expected empty string, got %q", gotS)
	}

	type user struct {
		ID int
	}
	var pu *user
	gotU := util.Deref(pu)
	if !reflect.DeepEqual(gotU, user{}) {
		t.Errorf("expected zero value %#v, got %#v", user{}, gotU)
	}
}
