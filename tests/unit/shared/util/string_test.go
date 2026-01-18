// Package util_test は文字列ユーティリティ
// (internal/shared/util/string.go) のユニットテストを提供する。
package util_test

import (
	"testing"

	"resume/internal/shared/util"
)

//
// ---------- OptPtr / NullableString / DerefString / SafeStr ----------
//

// TestOptPtr は空文字列の場合は nil、それ以外はポインタを返すことを確認する。
func TestOptPtr(t *testing.T) {
	t.Parallel()

	if got := util.OptPtr(""); got != nil {
		t.Errorf("expected nil for empty string, got=%v", got)
	}

	s := "value"
	got := util.OptPtr(s)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != s {
		t.Errorf("expected %q, got=%q", s, *got)
	}
}

// TestNullableString は空文字列の場合は nil、それ以外はポインタを返すことを確認する。
func TestNullableString(t *testing.T) {
	t.Parallel()

	if got := util.NullableString(""); got != nil {
		t.Errorf("expected nil for empty string, got=%v", got)
	}

	s := "hello"
	got := util.NullableString(s)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != s {
		t.Errorf("expected %q, got=%q", s, *got)
	}
}

// TestDerefString は nil ポインタの場合は空文字、非 nil の場合は値を返すことを確認する。
func TestDerefString(t *testing.T) {
	t.Parallel()

	if got := util.DerefString(nil); got != "" {
		t.Errorf("expected empty string for nil, got=%q", got)
	}

	s := "world"
	if got := util.DerefString(&s); got != s {
		t.Errorf("expected %q, got=%q", s, got)
	}
}

// TestSafeStr は入力をそのまま返すことを確認する。
func TestSafeStr(t *testing.T) {
	t.Parallel()

	if got := util.SafeStr("abc"); got != "abc" {
		t.Errorf("expected abc, got=%q", got)
	}
	if got := util.SafeStr(""); got != "" {
		t.Errorf("expected empty, got=%q", got)
	}
}

//
// ---------- FallbackStr ----------
//

// TestFallbackStr は最初の非空文字列を返すことを確認する。
func TestFallbackStr(t *testing.T) {
	t.Parallel()

	if got := util.FallbackStr(); got != "" {
		t.Errorf("expected empty for no args, got=%q", got)
	}

	if got := util.FallbackStr("", "", "first", "second"); got != "first" {
		t.Errorf("expected first, got=%q", got)
	}

	if got := util.FallbackStr("", "only"); got != "only" {
		t.Errorf("expected only, got=%q", got)
	}

	if got := util.FallbackStr("", ""); got != "" {
		t.Errorf("expected empty when all empty, got=%q", got)
	}
}

//
// ---------- ToSnake / ToKebab ----------
//

// TestToSnake_BasicCases はさまざまな表記を snake_case に変換できることを確認する。
func TestToSnake_BasicCases(t *testing.T) {
	t.Parallel()

	type tc struct {
		in, want string
	}
	tests := []tc{
		{"", ""},
		{"user", "user"},
		{"userName", "user_name"},
		{"UserName", "user_name"},
		{"user_name", "user_name"},
		{"user-name", "user_name"},
		{"HTTPServer", "http_server"},
		{"HTTPServerID", "http_server_id"},
		{"JP_PREF", "jp_pref"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := util.ToSnake(tt.in)
			if got != tt.want {
				t.Errorf("ToSnake(%q) expected %q, got %q", tt.in, tt.want, got)
			}
		})
	}
}

// TestToKebab_BasicCases はさまざまな表記を kebab-case に変換できることを確認する。
func TestToKebab_BasicCases(t *testing.T) {
	t.Parallel()

	type tc struct {
		in, want string
	}
	tests := []tc{
		{"", ""},
		{"user", "user"},
		{"userName", "user-name"},
		{"UserName", "user-name"},
		{"user_name", "user-name"},
		{"user-name", "user-name"},
		{"HTTPServer", "http-server"},
		{"HTTPServerID", "http-server-id"},
		{"JP_PREF", "jp-pref"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := util.ToKebab(tt.in)
			if got != tt.want {
				t.Errorf("ToKebab(%q) expected %q, got %q", tt.in, tt.want, got)
			}
		})
	}
}

//
// ---------- ToPascal ----------
//

// TestToPascal_BasicCases はさまざまな表記を PascalCase に変換できることを確認する。
func TestToPascal_BasicCases(t *testing.T) {
	t.Parallel()

	type tc struct {
		in, want string
	}
	tests := []tc{
		{"", ""},
		{"user", "User"},
		{"user_name", "UserName"},
		{"user-name", "UserName"},
		{"UserName", "UserName"},
		{"HTTPServer", "HttpServer"},
		{"jp_pref", "JpPref"},
		{"JP_PREF", "JpPref"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := util.ToPascal(tt.in)
			if got != tt.want {
				t.Errorf("ToPascal(%q) expected %q, got %q", tt.in, tt.want, got)
			}
		})
	}
}

//
// ---------- ToCamel ----------
//

// TestToCamel_BasicCases は snake/kebab/Pascal/UPPER/camel を lowerCamelCase に変換できることを確認する。
func TestToCamel_BasicCases(t *testing.T) {
	t.Parallel()

	type tc struct {
		in, want string
	}
	tests := []tc{
		{"", ""},
		{"jp_pref", "jpPref"},
		{"postal-code", "postalCode"},
		{"PostalCode", "postalCode"},
		{"JP_PREF", "jpPref"},
		{"jpPref", "jpPref"},
		{"user_name", "userName"},
		{"UserName", "userName"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := util.ToCamel(tt.in)
			if got != tt.want {
				t.Errorf("ToCamel(%q) expected %q, got %q", tt.in, tt.want, got)
			}
		})
	}
}

//
// ---------- Humanize ----------
//

// TestHumanize_BasicCases はさまざまな表記を人間可読な表記に変換できることを確認する。
func TestHumanize_BasicCases(t *testing.T) {
	t.Parallel()

	type tc struct {
		in, want string
	}
	tests := []tc{
		{"", ""},
		{"postalCode", "Postal code"},
		{"address_line1", "Address line1"},
		{"HTTPServerID", "Http server id"},
		{"jp-pref", "Jp pref"},
		{"JP_PREF", "Jp pref"},
		{"  multiple__spaces___here  ", "Multiple spaces here"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got := util.Humanize(tt.in)
			if got != tt.want {
				t.Errorf("Humanize(%q) expected %q, got %q", tt.in, tt.want, got)
			}
		})
	}
}

//
// ---------- Fmt ----------
//

// TestFmt は fmt.Sprintf 相当の文字列生成ができることを確認する。
func TestFmt(t *testing.T) {
	t.Parallel()

	got := util.Fmt("ui.profile.age_group.%d", 20)
	want := "ui.profile.age_group.20"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
