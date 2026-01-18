// Package entity_test は domain/entity のユニットテストを提供します。
// このファイルでは UserEducation エンティティに関するテストを行います。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

func TestUserEducation_TableName(t *testing.T) {
	t.Parallel()

	want := "user_educations"
	got := (entity.UserEducation{}).TableName()

	if got != want {
		t.Errorf("UserEducation.TableName() = %s, want %s", got, want)
	}
}
