// UserProfile エンティティに関するユニットテストです。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestUserProfile_TableName は UserProfile のテーブル名が正しいことを確認します。
func TestUserProfile_TableName(t *testing.T) {
	t.Parallel()

	want := "user_profiles"
	got := (entity.UserProfile{}).TableName()

	if got != want {
		t.Errorf("UserProfile.TableName() = %s, want %s", got, want)
	}
}
