// Package entity_test は domain/entity のユニットテストを提供します。
// このファイルでは AddressPurpose エンティティに関するテストを行います。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestAddressPurpose_TableName は TableName() の返却値が
// GORM で期待されるテーブル名と一致することを検証します。
func TestAddressPurpose_TableName(t *testing.T) {
	t.Parallel()

	want := "address_purposes"
	got := (entity.AddressPurpose{}).TableName()

	if got != want {
		t.Errorf("AddressPurpose.TableName() = %s, want %s", got, want)
	}
}

// TestAddressPurpose_Constants は、AddressPurpose に定義されている
// ID 定数および Code 定数が最低限有効な値を持つことを確認するテストです。
//
// 本テストで保証するのは、次の「破壊的変更を防ぐための契約」のみです。
// - AddressPurpose の各 ID 定数が 0 ではないこと
// - AddressPurpose の各 Code 定数が空文字ではないこと
//
// 各用途の意味的な正当性や表示名、並び順などは
// データベースのマスタ定義および seed データの責務であり、
// Go のユニットテストでは検証対象としません。
//
// このテストの目的は、定数の追加・削除・リネーム時に発生しがちな
// 定義漏れやゼロ値代入といった不具合を早期に検出するための
// 安全ネットを提供することです。
func TestAddressPurpose_Constants(t *testing.T) {
	tests := []struct {
		id   uint64
		code string
	}{
		{entity.AddressPurposeIDHome, entity.AddressPurposeCodeHome},
		{entity.AddressPurposeIDContact, entity.AddressPurposeCodeContact},
		{entity.AddressPurposeIDOther, entity.AddressPurposeCodeOther},
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
