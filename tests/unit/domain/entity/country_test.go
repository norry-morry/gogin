// Country エンティティに関するユニットテストです。
package entity_test

import (
	"testing"

	"resume/internal/domain/entity"
)

// TestCountry_TableName は Country のテーブル名が正しいことを確認します。
func TestCountry_TableName(t *testing.T) {
	t.Parallel()

	want := "countries"
	got := (entity.Country{}).TableName()

	if got != want {
		t.Errorf("Country.TableName() = %s, want %s", got, want)
	}
}

// TestCountry_DisplayCode は DisplayCode() が Code をそのまま返すことを確認します。
func TestCountry_DisplayCode(t *testing.T) {
	t.Parallel()

	c := entity.Country{
		Code:        "JP",
		IsSupported: true,
	}

	if got := c.DisplayCode(); got != "JP" {
		t.Errorf("DisplayCode() = %s, want JP", got)
	}
}

// TestCountry_IsActive は IsActive() が IsSupported を返すことを確認します。
func TestCountry_IsActive(t *testing.T) {
	t.Parallel()

	active := entity.Country{
		Code:        "JP",
		IsSupported: true,
	}
	inactive := entity.Country{
		Code:        "US",
		IsSupported: false,
	}

	if !active.IsActive() {
		t.Errorf("expected active country to be IsActive=true, got false")
	}
	if inactive.IsActive() {
		t.Errorf("expected inactive country to be IsActive=false, got true")
	}
}
