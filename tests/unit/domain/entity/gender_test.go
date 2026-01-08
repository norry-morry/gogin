// Gender エンティティに関するユニットテストです。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestGender_TableName は Gender のテーブル名が正しいことを確認します。
func TestGender_TableName(t *testing.T) {
	t.Parallel()

	want := "genders"
	got := (entity.Gender{}).TableName()

	if got != want {
		t.Errorf("Gender.TableName() = %s, want %s", got, want)
	}
}

// TestGender_Constants は、Gender に定義されている
// 定数値（ID / Code）が最低限有効な値を持つことを確認するテストです。
//
// 本テストで検証する内容は以下に限定されます。
// - ID 定数がゼロ値（0）ではないこと
// - Code 定数が空文字ではないこと
//
// Gender の表示名、並び順、DB 上の整合性、
// あるいは「男女以外を含めるべきか」といった仕様判断は
// ドメイン設計・マスタデータの責務であり、
// Go のユニットテストでは扱いません。
//
// このテストは、定数の追加・削除・リネーム時に
// ゼロ値や空文字が誤って混入することを防ぐための
// リグレッションテストとして位置付けられます。
func TestGender_Constants(t *testing.T) {
	tests := []struct {
		name string
		id   uint8
		code string
	}{
		{
			name: "Male",
			id:   entity.GenderMaleID,
			code: entity.GenderMaleCode,
		},
		{
			name: "Female",
			id:   entity.GenderFemaleID,
			code: entity.GenderFemaleCode,
		},
		{
			name: "Other",
			id:   entity.GenderOtherID,
			code: entity.GenderOtherCode,
		},
	}

	for _, tt := range tests {
		tt := tt // parallel safety
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.id == 0 {
				t.Errorf("Gender ID must not be zero")
			}
			if tt.code == "" {
				t.Errorf("Gender Code must not be empty")
			}
		})
	}
}
