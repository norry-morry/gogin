// tests/unit/adapter/http/presenter/user_profile_presenter_test.go
package presenter_test

import (
	"testing"
	"time"

	"resume/internal/adapter/http/presenter"
	"resume/internal/usecase/profile"
)

func TestPresentUserProfile(t *testing.T) {
	b := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	age := 35
	genderID := uint8(1)

	out := profile.UserProfileOutput{
		UserID:    1,
		BirthDate: &b,
		Age:       &age,
		GenderID:  &genderID,
	}

	got := presenter.PresentUserProfile(out)

	if got.BirthDate == nil {
		t.Fatalf("BirthDate should be formatted string")
	}

	if got.AgeGroup == nil || *got.AgeGroup != 30 {
		t.Errorf("AgeGroup mismatch: want 30")
	}
}
