// Package entity_test は domain/entity のユニットテストを提供します。
// このファイルでは EducationStatus エンティティに関するテストを行います。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestEducationStatus_TableName は TableName() の返却値が
// GORM で期待されるテーブル名と一致することを検証します。
func TestEducationStatus_TableName(t *testing.T) {
	t.Parallel()

	want := "education_statuses"
	got := (entity.EducationStatus{}).TableName()

	if got != want {
		t.Errorf("EducationStatus.TableName() = %s, want %s", got, want)
	}
}

// TestEducationStatus_Constants は、EducationStatus に定義されている
// ID 定数および Code 定数が「最低限の契約」を満たしていることを確認するテストです。
//
// 本テストでは、以下の点のみを検証対象とします。
// - 各ステータス ID が 0 ではないこと
// - 各ステータス Code が空文字ではないこと
//
// 値そのもの（数値や文字列の内容）が正しいかどうかは、
// マイグレーションや初期データ（seed）の責務であり、
// Go のユニットテストでは検証しません。
//
// このテストの目的は、定数の追加・変更時に
// 「未定義」「ゼロ値」「空文字」といった致命的な破壊を
// 早期に検知するための安全ネットを提供することです。
func TestEducationStatus_Constants(t *testing.T) {
	tests := []struct {
		id   uint64
		code string
	}{
		{entity.EducationStatusEnrolled, entity.EducationStatusCodeEnrolled},
		{entity.EducationStatusGraduated, entity.EducationStatusCodeGraduated},
		{entity.EducationStatusLeaveOfAbsence, entity.EducationStatusCodeLeaveOfAbsence},
	}

	for _, tt := range tests {
		if tt.id == 0 {
			t.Errorf("id must not be zero for code = %s", tt.code)
		}
		if tt.code == "" {
			t.Errorf("id must not be empty for id = %d", tt.id)
		}
	}
}
