// Package entity_test は domain/entity のユニットテストを提供します。
// このファイルでは DegreeType エンティティに関するテストを行います。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestDegreeType_TableName は TableName() の返却値が
// GORM で期待されるテーブル名と一致することを検証します。
func TestDegreeType_TableName(t *testing.T) {
	t.Parallel()

	want := "degree_types"
	got := (entity.DegreeType{}).TableName()

	if got != want {
		t.Errorf("DegreeType.TableName() = %s, want %s", got, want)
	}
}

// TestDegreeType_Constants は、DegreeType に定義されている
// ID 定数および Code 定数が最低限有効な値を持つことを確認するテストです。
//
// 本テストでは、以下の契約のみを保証します。
// - DegreeType の ID 定数が 0 ではないこと
// - DegreeType の Code 定数が空文字ではないこと
//
// 各定数の具体的な値や並び順の正当性については、
// データベース設計および seed データで管理されるべきであり、
// 本テストでは扱いません。
//
// このテストは、将来的な定数追加やリファクタリング時に
// 定義漏れや初期化ミスを検出するための軽量なユニットテストです。
func TestDegreeType_constants(t *testing.T) {
	tests := []struct {
		id   uint64
		code string
	}{
		{entity.DegreeTypeHighSchool, entity.DegreeTypeCodeHighSchool},
		{entity.DegreeTypeVocational, entity.DegreeTypeCodeVocational},
		{entity.DegreeTypeOther, entity.DegreeTypeCodeOther},
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
