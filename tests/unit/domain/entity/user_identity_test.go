// UserIdentity エンティティに関するユニットテストです。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestUserIdentity_TableName は UserIdentity のテーブル名が正しいことを確認します。
func TestUserIdentity_TableName(t *testing.T) {
	t.Parallel()

	want := "user_identities"
	got := (entity.UserIdentity{}).TableName()

	if got != want {
		t.Errorf("UserIdentity.TableName() = %s, want %s", got, want)
	}
}
