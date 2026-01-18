// Package util はアプリケーション全体で共有される汎用的なユーティリティを集めたパッケージです。
// 本ファイル (date.go) では日付関連の関数（年齢・年代計算など誕生日からの年齢・年代計算など、ビュー層での派生情報生成）を提供します。
package util

import (
	"time"
)

// CalcAge は誕生日から現在の満年齢を返します。
// birth が nil または未来日の場合は 0 を返します。
func CalcAge(birth *time.Time) int {
	if birth == nil {
		return 0
	}
	now := time.Now()
	y, m, d := now.Date()
	by, bm, bd := birth.Date()
	age := y - by
	if (m < bm) || (m == bm && d < bd) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

// CalcAgeGroup は年齢から年代（10の倍数）を整数で返します。
// 例: 25 → 20, 38 → 30, 61 → 60。
// 年齢が0以下の場合は 0 を返します。
func CalcAgeGroup(age int) int {
	if age <= 0 {
		return 0
	}
	if age < 10 {
		return 0 // 10歳未満扱い
	}
	return (age / 10) * 10
}

// FormatDatePtr は *time.Time を ISO8601 (YYYY-MM-DD) 形式の *string に変換します。
// nil の場合は nil を返します。
func FormatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

// FormatDatePtrWithLayout は指定レイアウトで *time.Time を *string に変換します。
func FormatDatePtrWithLayout(t *time.Time, layout string) *string {
	if t == nil {
		return nil
	}
	s := t.Format(layout)
	return &s
}

// ParseYearMonth は "2006-01" 形式の文字列ポインタをパースして *time.Time を返します。
// 入力が nil または空文字の場合は nil を返します。
// パースエラー時は error を返します。
func ParseYearMonth(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ParseYearMonthVal は "2006-01" 形式の文字列ポインタをパースして time.Time (値) を返します。
// 入力が nil/空文字/エラーの場合はゼロ値を返します（必須項目用）。
// エラーハンドリングを厳密にしたい場合は上記 ParseYearMonth を使って呼び出し元で処理してください。
func ParseYearMonthVal(s *string) (time.Time, error) {
	if s == nil || *s == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01", *s)
}
