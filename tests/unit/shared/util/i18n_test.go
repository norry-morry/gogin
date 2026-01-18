// Package util_test は i18n 関連ユーティリティ (internal/shared/util/i18n.go)
// の公開関数に対するユニットテストを提供する。
package util_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"resume/internal/shared/util"
)

//
// ---------- ReadYAMLFile ----------
//

// TestReadYAMLFile_Success は ReadYAMLFile が YAML を map にパースできることを確認する。
func TestReadYAMLFile_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.yaml")

	content := `
ui:
  page:
    home:
      title: Home
count: 1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp yaml: %v", err)
	}

	got, err := util.ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile returned error: %v", err)
	}

	ui, ok := got["ui"].(map[string]any)
	if !ok {
		t.Fatalf("ui is not map[string]any: %#v", got["ui"])
	}
	pageVal, ok := ui["page"].(map[string]any)
	if !ok {
		t.Fatalf("ui[\"page\"] is not map[string]any: %#v", ui["page"])
	}
	page := pageVal

	homeVal, ok := page["home"].(map[string]any)
	if !ok {
		t.Fatalf("page[\"home\"] is not map[string]any: %#v", page["home"])
	}
	home := homeVal

	if home["title"] != "Home" {
		t.Errorf("expected title=Home, got=%v", home["title"])
	}

	if got["count"] != 1 {
		t.Errorf("expected count=1, got=%v", got["count"])
	}
}

// TestReadYAMLFile_NotFound は存在しないパスを渡した場合にエラーを返すことを確認する。
func TestReadYAMLFile_NotFound(t *testing.T) {
	t.Parallel()

	got, err := util.ReadYAMLFile("no_such_file.yaml")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if got != nil {
		t.Errorf("expected nil map on error, got=%v", got)
	}
}

//
// ---------- WrapWithPath ----------
//

// TestWrapWithPath は prefix と map を組み合わせてネスト構造を生成できることを確認する。
func TestWrapWithPath(t *testing.T) {
	t.Parallel()

	in := map[string]any{"title": "Home"}

	t.Run("プレフィックスなしならそのまま返す", func(t *testing.T) {
		got := util.WrapWithPath(nil, in)
		if !reflect.DeepEqual(got, in) {
			t.Errorf("expected %v, got %v", in, got)
		}
	})

	t.Run("単一要素のプレフィックス", func(t *testing.T) {
		got := util.WrapWithPath([]string{"home"}, in)
		want := map[string]any{
			"home": map[string]any{
				"title": "Home",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("複数要素のプレフィックス", func(t *testing.T) {
		got := util.WrapWithPath([]string{"ui", "page", "home"}, in)
		want := map[string]any{
			"ui": map[string]any{
				"page": map[string]any{
					"home": map[string]any{
						"title": "Home",
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

//
// ---------- DeepMerge ----------
//

// TestDeepMerge_MergeNestedMaps はネストした map 同士を再帰的にマージできることを確認する。
func TestDeepMerge_MergeNestedMaps(t *testing.T) {
	t.Parallel()

	dst := map[string]any{
		"ui": map[string]any{
			"page": map[string]any{
				"home": map[string]any{
					"title": "Home",
				},
			},
		},
	}
	src := map[string]any{
		"ui": map[string]any{
			"page": map[string]any{
				"home": map[string]any{
					"description": "desc",
				},
			},
		},
		"common": "value",
	}

	if err := util.DeepMerge(dst, src); err != nil {
		t.Fatalf("DeepMerge returned error: %v", err)
	}

	want := map[string]any{
		"ui": map[string]any{
			"page": map[string]any{
				"home": map[string]any{
					"title":       "Home",
					"description": "desc",
				},
			},
		},
		"common": "value",
	}

	if !reflect.DeepEqual(dst, want) {
		t.Errorf("expected %v, got %v", want, dst)
	}
}

// TestDeepMerge_OverrideNonMap は src が非 map の場合に上書きされることを確認する。
func TestDeepMerge_OverrideNonMap(t *testing.T) {
	t.Parallel()

	dst := map[string]any{
		"ui": map[string]any{
			"title": "Home",
		},
	}
	src := map[string]any{
		"ui": "string value",
	}

	if err := util.DeepMerge(dst, src); err != nil {
		t.Fatalf("DeepMerge returned error: %v", err)
	}

	if dst["ui"] != "string value" {
		t.Errorf("expected ui to be overridden, got=%v", dst["ui"])
	}
}

//
// ---------- Flatten ----------
//

// TestFlatten_Nested はネスト構造が "a.b.c" 形式にフラット化されることを確認する。
func TestFlatten_Nested(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"ui": map[string]any{
			"page": map[string]any{
				"home": map[string]any{
					"title": "Home",
				},
			},
		},
	}
	out := make(map[string]any)

	util.Flatten("", in, out)

	want := map[string]any{
		"ui.page.home.title": "Home",
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("expected %v, got %v", want, out)
	}
}

// TestFlatten_MapAnyKeys は map[any]any が混在しても正しくフラット化されることを確認する。
func TestFlatten_MapAnyKeys(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"root": map[any]any{
			"child": "value",
		},
	}
	out := make(map[string]any)

	util.Flatten("", in, out)

	want := map[string]any{
		"root.child": "value",
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("expected %v, got %v", want, out)
	}
}

//
// ---------- DetectLocale ----------
//

// TestDetectLocale_Ja はパスに ja を含む場合に ja が返ることを確認する。
func TestDetectLocale_Ja(t *testing.T) {
	t.Parallel()

	path := "internal/shared/i18n/ja/master/address_purpose.yaml"
	got := util.DetectLocale(path)

	if got != "ja" {
		t.Errorf("expected ja, got=%s", got)
	}
}

// TestDetectLocale_En はパスに en を含む場合に en が返ることを確認する。
func TestDetectLocale_En(t *testing.T) {
	t.Parallel()

	path := "internal/shared/i18n/en/master/country.yaml"
	got := util.DetectLocale(path)

	if got != "en" {
		t.Errorf("expected en, got=%s", got)
	}
}

// TestDetectLocale_WindowsStylePath は Windows 風のパスでも en/ja を検出できることを確認する。
func TestDetectLocale_WindowsStylePath(t *testing.T) {
	t.Parallel()

	path := `internal\shared\i18n\locales\en\master\country.yaml`
	got := util.DetectLocale(path)

	if got != "en" {
		t.Errorf("expected en, got=%s", got)
	}
}

// TestDetectLocale_DefaultJa は候補が含まれない場合は ja が返ることを確認する。
func TestDetectLocale_DefaultJa(t *testing.T) {
	t.Parallel()

	path := "internal/shared/i18n/master/country.yaml"
	got := util.DetectLocale(path)

	if got != "ja" {
		t.Errorf("expected default ja, got=%s", got)
	}
}
