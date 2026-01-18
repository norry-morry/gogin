// Package util_test はマスキング関連のグローバルユーティリティ
// (internal/shared/util/masker_global.go) のユニットテストを提供する。
package util_test

import (
	"testing"

	"resume/internal/shared/util"
)

// resetGlobalMaskerForTest はテスト用にグローバル Masker を Noop に戻すヘルパー。
func resetGlobalMaskerForTest() {
	util.SetGlobalMasker(util.NoopMasker{})
}

//
// ---------- GlobalMasker / SetGlobalMasker の基本動作 ----------
//

// TestGlobalMasker_Default は初期状態で GlobalMasker が NoopMasker を返すことを確認する。
func TestGlobalMasker_Default(t *testing.T) {
	resetGlobalMaskerForTest()

	m := util.GlobalMasker()
	if m == nil {
		t.Fatal("expected non-nil Masker")
	}

	if _, ok := m.(util.NoopMasker); !ok {
		t.Errorf("expected default masker to be NoopMasker, got %T", m)
	}
}

// testMasker は SetGlobalMasker 用のカスタム Masker 型。
// NoopMasker を埋め込み、既存のインターフェース実装をそのまま利用する。
type testMasker struct {
	util.NoopMasker
}

// TestSetGlobalMasker_Overrides は SetGlobalMasker でグローバル Masker が差し替わることを確認する。
func TestSetGlobalMasker_Overrides(t *testing.T) {
	resetGlobalMaskerForTest()
	defer resetGlobalMaskerForTest()

	custom := testMasker{}
	util.SetGlobalMasker(custom)

	got := util.GlobalMasker()
	if _, ok := got.(testMasker); !ok {
		t.Errorf("expected GlobalMasker to return testMasker, got %T", got)
	}
}

// TestSetGlobalMasker_NilIsIgnored は nil を渡した場合にグローバル Masker が変更されないことを確認する。
func TestSetGlobalMasker_NilIsIgnored(t *testing.T) {
	resetGlobalMaskerForTest()
	defer resetGlobalMaskerForTest()

	before := util.GlobalMasker()
	if before == nil {
		t.Fatal("expected non-nil Masker before SetGlobalMasker(nil)")
	}

	// nil を渡しても変更されないはず
	util.SetGlobalMasker(nil)

	after := util.GlobalMasker()
	if after == nil {
		t.Fatal("expected non-nil Masker after SetGlobalMasker(nil)")
	}

	if before != after {
		t.Errorf("expected global masker not to change when passing nil, before=%T after=%T", before, after)
	}
}
