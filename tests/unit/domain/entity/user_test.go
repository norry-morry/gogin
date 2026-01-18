// User エンティティに関するユニットテストです。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestUser_TableName は User のテーブル名が正しいことを確認します。
func TestUser_TableName(t *testing.T) {
	t.Parallel()

	want := "users"
	got := (entity.User{}).TableName()

	if got != want {
		t.Errorf("User.TableName() = %s, want %s", got, want)
	}
}
