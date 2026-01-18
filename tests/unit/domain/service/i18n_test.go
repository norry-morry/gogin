// Package service_test は i18n ドメインサービスのユニットテストを提供します。
package service_test

import (
	"testing"

	svc "resume/internal/domain/service"
	vo "resume/internal/domain/valueobject/i18n"
)

//
// ---------- 0個の場合 ----------
//

// TestMergeBundles_EmptyInput は bundles が 0 件の場合に空バンドルが返ることを確認します。
func TestMergeBundles_EmptyInput(t *testing.T) {
	t.Parallel()

	b := svc.MergeBundles()

	if b.Locale().Code() != "" {
		t.Errorf("expected empty locale, got %s", b.Locale().Code())
	}

	if len(b.Messages()) != 0 {
		t.Errorf("expected empty messages, got=%v", b.Messages())
	}
}

//
// ---------- 1個の場合 ----------
//

// TestMergeBundles_SingleBundle は 1 つだけの bundle をそのまま返すことを確認します。
func TestMergeBundles_SingleBundle(t *testing.T) {
	t.Parallel()

	locale := vo.NewLocale("ja")
	msgs := map[string]string{
		"a": "A",
		"b": "B",
	}
	b1 := vo.NewBundle(locale, msgs)

	out := svc.MergeBundles(b1)

	if !out.Locale().Equal(locale) {
		t.Errorf("expected locale=ja, got %s", out.Locale().Code())
	}

	if outV, _ := out.Get(vo.NewKey("a")); outV != "A" {
		t.Errorf("expected a=A, got=%s", outV)
	}
	if outV, _ := out.Get(vo.NewKey("b")); outV != "B" {
		t.Errorf("expected b=B, got=%s", outV)
	}
}

//
// ---------- 複数マージ（後勝ち） ----------
//

// TestMergeBundles_OverrideLaterBundles は 後の bundle が前の値を上書きすることを確認します。
func TestMergeBundles_OverrideLaterBundles(t *testing.T) {
	t.Parallel()

	locale := vo.NewLocale("ja")

	b1 := vo.NewBundle(locale, map[string]string{
		"a": "A1",
		"b": "B1",
	})

	b2 := vo.NewBundle(locale, map[string]string{
		"b": "B2", // 上書き
		"c": "C2",
	})

	b3 := vo.NewBundle(locale, map[string]string{
		"a": "A3", // 最終上書き
	})

	out := svc.MergeBundles(b1, b2, b3)

	tests := map[string]string{
		"a": "A3", // b3 が勝ち
		"b": "B2", // b2 が勝ち
		"c": "C2", // b2 から
	}

	for k, expect := range tests {
		got, ok := out.Get(vo.NewKey(k))
		if !ok {
			t.Fatalf("key %s missing", k)
		}
		if got != expect {
			t.Errorf("key %s: expected %s, got %s", k, expect, got)
		}
	}
}

//
// ---------- Locale の扱い ----------
//

// TestMergeBundles_LocaleIsTakenFromFirst は Locale が最初の bundle のものになることを確認します。
func TestMergeBundles_LocaleIsTakenFromFirst(t *testing.T) {
	t.Parallel()

	ja := vo.NewLocale("ja")
	en := vo.NewLocale("en")

	b1 := vo.NewBundle(ja, map[string]string{"a": "A"})
	b2 := vo.NewBundle(en, map[string]string{"b": "B"})

	out := svc.MergeBundles(b1, b2)

	if !out.Locale().Equal(ja) {
		t.Errorf("expected locale from first bundle (ja), got %s", out.Locale().Code())
	}
}

//
// ---------- map 破壊防止（イミュータビリティ） ----------
//

// TestMergeBundles_Immutability は MergeBundles が内部 map のコピーを使っていることを確認します。
func TestMergeBundles_Immutability(t *testing.T) {
	t.Parallel()

	locale := vo.NewLocale("ja")

	src1 := map[string]string{"a": "A"}
	src2 := map[string]string{"b": "B"}

	b1 := vo.NewBundle(locale, src1)
	b2 := vo.NewBundle(locale, src2)

	out := svc.MergeBundles(b1, b2)

	// 元 map を書き換え
	src1["a"] = "A_CHANGED"
	src2["b"] = "B_CHANGED"

	// Bundle の中は影響を受けないはず
	if v, _ := out.Get(vo.NewKey("a")); v != "A" {
		t.Errorf("expected A, got=%s", v)
	}
	if v, _ := out.Get(vo.NewKey("b")); v != "B" {
		t.Errorf("expected B, got=%s", v)
	}

	// Messages() の変更が out に影響しないことも確認
	msgs := out.Messages()
	msgs["a"] = "PATCHED"

	if v, _ := out.Get(vo.NewKey("a")); v != "A" {
		t.Errorf("expected immutability, got changed value=%s", v)
	}
}
