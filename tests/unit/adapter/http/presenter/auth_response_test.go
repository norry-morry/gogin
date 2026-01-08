// tests/unit/adapter/http/presenter/auth_response_test.go

package presenter_test

import (
	"testing"

	"resume/internal/adapter/http/presenter"
	"resume/internal/domain/entity"
)

func TestFromDomainUserForAuth(t *testing.T) {
	token := "abc123"

	u := &entity.User{
		UID:         "uid_001",
		DisplayName: ptr("Taro"),
		Email:       ptr("test@example.com"),
		PhotoURL:    ptr("http://example.com/photo.jpg"),
	}

	got := presenter.FromDomainUserForAuth(u, token)

	if got.UID != u.UID {
		t.Errorf("UID mismatch: want=%s, got=%s", u.UID, got.UID)
	}
	if *got.DisplayName != *u.DisplayName {
		t.Errorf("DisplayName mismatch")
	}
	if got.Token != token {
		t.Errorf("Token mismatch")
	}
}

func ptr[T any](v T) *T { return &v }
