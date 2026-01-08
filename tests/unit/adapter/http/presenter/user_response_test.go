// tests/unit/adapter/http/presenter/user_response_test.go

package presenter_test

import (
	"testing"
	"time"

	"resume/internal/adapter/http/presenter"
	"resume/internal/domain/entity"
)

func TestToUserResponse(t *testing.T) {
	now := time.Now()

	u := &entity.User{
		ID:            1,
		UID:           "uid01",
		Email:         ptr("a@b.com"),
		DisplayName:   ptr("Taro"),
		PhotoURL:      ptr("url"),
		Disabled:      false,
		EmailVerified: true,
		LastLoginAt:   &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	got := presenter.ToUserResponse(u)
	if *got.Email != *u.Email {
		t.Errorf("Email mismatch")
	}
	if got.CreatedAt == "" {
		t.Errorf("CreatedAt should be formatted")
	}
}
