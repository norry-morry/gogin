// Package util_test は util パッケージの公開関数に対するユニットテストを提供する。
// 本ファイルでは date.go に定義された日付ユーティリティ関数（CalcAge、CalcAgeGroup、FormatDatePtr 等）の
// 振る舞いをテストする。
package util_test

import (
	"testing"
	"time"

	"resume/internal/shared/util"
)

//
// ---------- CalcAge ----------
//

// TestCalcAge は CalcAge の動作（nil, 未来日, 誕生日計算, 誕生日未到来）を確認する。
func TestCalcAge(t *testing.T) {
	t.Run("nil を渡した場合は 0 を返す", func(t *testing.T) {
		got := util.CalcAge(nil)
		if got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})

	t.Run("未来日を渡した場合は 0 を返す", func(t *testing.T) {
		future := time.Now().AddDate(1, 0, 0)
		got := util.CalcAge(&future)
		if got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})

	t.Run("誕生日が今日と同じ場合は年齢がそのまま計算される", func(t *testing.T) {
		now := time.Now()
		birth := now.AddDate(-30, 0, 0)

		got := util.CalcAge(&birth)
		if got != 30 {
			t.Errorf("expected 30, got %d", got)
		}
	})

	t.Run("誕生日がまだ来ていない場合は年齢が1減算される", func(t *testing.T) {
		now := time.Now()

		birth := time.Date(
			now.Year()-20,
			now.Month()+1, // 来月 → 誕生日がまだ来ていない
			now.Day(),
			0, 0, 0, 0, time.UTC,
		)

		got := util.CalcAge(&birth)
		if got != 19 {
			t.Errorf("expected 19, got %d", got)
		}
	})
}

//
// ---------- CalcAgeGroup ----------
//

// TestCalcAgeGroup は CalcAgeGroup の年齢→年代変換の動作をテストする。
func TestCalcAgeGroup(t *testing.T) {
	tests := []struct {
		age      int
		expected int
	}{
		{0, 0},
		{-5, 0},
		{5, 0},
		{10, 10},
		{15, 10},
		{25, 20},
		{38, 30},
		{61, 60},
	}

	for _, tt := range tests {
		tt := tt // ローカルコピー（parallel 対策）
		t.Run(
			"age group",
			func(t *testing.T) {
				got := util.CalcAgeGroup(tt.age)
				if got != tt.expected {
					t.Errorf("age=%d expected %d, got %d", tt.age, tt.expected, got)
				}
			},
		)
	}
}

//
// ---------- FormatDatePtr ----------
//

// TestFormatDatePtr は FormatDatePtr の nil / 日付フォーマット動作を確認する。
func TestFormatDatePtr(t *testing.T) {
	t.Run("nil を渡した場合は nil を返す", func(t *testing.T) {
		got := util.FormatDatePtr(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("*time.Time を YYYY-MM-DD にフォーマットして返す", func(t *testing.T) {
		d := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		got := util.FormatDatePtr(&d)
		expected := "2024-12-31"

		if got == nil || *got != expected {
			t.Errorf("expected %s, got %v", expected, got)
		}
	})
}

//
// ---------- FormatDatePtrWithLayout ----------
//

// TestFormatDatePtrWithLayout は指定レイアウトでのフォーマット変換動作を確認する。
func TestFormatDatePtrWithLayout(t *testing.T) {
	t.Run("nil を渡した場合は nil を返す", func(t *testing.T) {
		got := util.FormatDatePtrWithLayout(nil, "2006/01/02")
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("指定レイアウトでフォーマットされる", func(t *testing.T) {
		d := time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)
		got := util.FormatDatePtrWithLayout(&d, "2006/01/02")
		expected := "2025/01/13"

		if got == nil || *got != expected {
			t.Errorf("expected %s, got %v", expected, got)
		}
	})
}
