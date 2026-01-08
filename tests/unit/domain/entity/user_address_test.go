// Package entity_test は internal/domain/entity パッケージに定義された
// ドメインエンティティのユニットテストを提供します。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

func TestUserAddress_TableName(t *testing.T) {
	t.Parallel()

	var ua entity.UserAddress

	if got, want := ua.TableName(), "user_addresses"; got != want {
		t.Fatalf("TableName mismatch: want=%s, got=%s", want, got)
	}
}
