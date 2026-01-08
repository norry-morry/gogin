// tests/unit/adapter/http/presenter/address_response_test.go
package presenter_test

import (
	"testing"
	"time"

	"resume/internal/adapter/http/presenter"
	"resume/internal/domain/entity"
)

func TestToAddressResponse(t *testing.T) {
	now := time.Now()
	postal := "100-0001"
	a := &entity.UserAddress{
		ID:          10,
		CountryCode: "JP",
		PostalCode:  &postal,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	got := presenter.ToAddressResponse(a)
	if got.CountryCode != "JP" {
		t.Errorf("CountryCode mismatch")
	}
}
