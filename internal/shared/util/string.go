// Package util は、共通の軽量ユーティリティ関数群を提供します。
// string.go は文字列関連の便利関数をまとめています
package util

import (
	"fmt"
	"strings"
	"unicode"
)

// OptPtr は空文字列なら nil、非空なら *string を返します。
// JSONやDBにNULLを入れたいときなどに便利です。
func OptPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// SafeStr は、将来のnil対応などを見越した「安全な文字列取得」用のラッパ。
// 現状はそのまま返しますが、nilガード付きのCloneなどと統一しておくと保守性が上がります。
func SafeStr(s string) string {
	return s
}

// FallbackStr は、最初に空でない文字列を返します。
// すべて空の場合は空文字列("")を返します。
//
// 例:
//
//	name := util.FallbackStr(p.DisplayName, user.DisplayName, "unknown")
func FallbackStr(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

// ToCamelLower は ToCamel へ統合したため削除（呼び出し側は ToCamel に置換
// ToCamelLower は、与えられた文字列の先頭文字を小文字に変換して返します。
// 例: "UserName" → "userName"。
// 空文字列を渡した場合はそのまま空文字を返します。
//func ToCamelLower(s string) string {
//	if s == "" {
//		return s
//	}
//	return strings.ToLower(s[:1]) + s[1:]
//}

// normalizeToSnake は、入力が snake/kebab/Pascal/camel/UPPER でも
// いったん snake_case（小文字）へ正規化する内部関数。
func normalizeToSnake(s string) string {
	if s == "" {
		return s
	}
	// kebab → snake
	s = strings.ReplaceAll(s, "-", "_")
	// Pascal/camel → snake（大文字手前に "_" を挿入）
	var out []rune
	var prev rune
	for i, r := range s {
		if r == '_' {
			out = append(out, r)
			prev = r
			continue
		}
		// 前が 英小文字/数字、今が 英大文字 → 境界
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)) {
			out = append(out, '_')
		}
		// 連続大文字 → 次が小文字なら境界（HTTPServer → http_server）
		if i > 0 && unicode.IsUpper(r) && unicode.IsUpper(prev) {
			// 次を先読み
			if i+1 < len(s) {
				n := rune(s[i+1])
				if unicode.IsLower(n) {
					out = append(out, '_')
				}
			}
		}
		out = append(out, unicode.ToLower(r))
		prev = r
	}
	// 連続区切りの正規化（"__" 等）: Split 時に無視されるのでここではそのままでもOK
	return string(out)
}

// ToSnake は、入力が何であっても snake_case に変換します。
func ToSnake(s string) string {
	return normalizeToSnake(s)
}

// ToKebab は、入力が何であっても kebab-case に変換します。
func ToKebab(s string) string {
	if s == "" {
		return s
	}
	return strings.ReplaceAll(normalizeToSnake(s), "_", "-")
}

// ToPascal は、入力が何であっても PascalCase に変換します。
func ToPascal(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(normalizeToSnake(s), "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

// ToCamel は入力文字列を lowerCamelCase に変換します。
// snake_case, kebab-case, PascalCase, UPPER_CASE いずれにも対応。
// 例:
//
//	"jp_pref"        → "jpPref"
//	"postal-code"    → "postalCode"
//	"PostalCode"     → "postalCode"
//	"JP_PREF"        → "jpPref"
//	"jpPref"         → "jpPref" (そのまま)
func ToCamel(s string) string {
	if s == "" {
		return s
	}

	// 1) まず snake に正規化
	parts := strings.Split(normalizeToSnake(s), "_")

	for i := range parts {
		if parts[i] == "" {
			continue
		}
		if i > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		} else {
			parts[i] = strings.ToLower(parts[i])
		}
	}
	return strings.Join(parts, "")

}

// Humanize は表記ゆれ（snake/kebab/camel/Pascal/UPPER）を問わず
// 「先頭だけ大文字 + スペース区切り」の人間可読表記に変換します。
// 例:
//   - "postalCode"     → "Postal code"
//   - "address_line1"  → "Address line1"
//   - "HTTPServerID"   → "Http server id"
//   - "jp-pref"        → "Jp pref"
//   - "JP_PREF"        → "Jp pref"
func Humanize(s string) string {
	if s == "" {
		return s
	}
	// 1) まず snake_case に正規化（小文字化もされる）
	s = normalizeToSnake(s)

	// 2) "_" → " " に置換（連続区切りは自然に連続スペースになり得る）
	s = strings.ReplaceAll(s, "_", " ")

	// 3) 余分なスペースの整理（前後・連続）
	s = strings.TrimSpace(s)
	// 連続スペースを 1 個に
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' {
			if !prevSpace {
				_, _ = b.WriteRune(' ')
				prevSpace = true
			}
			continue
		}
		_, _ = b.WriteRune(r)
		prevSpace = false
	}
	s = b.String()

	if s == "" {
		return s
	}

	// 4) 先頭だけ大文字に（他はすでに小文字化済みなのでそのまま）
	//    例: "http server id" → "Http server id"
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Fmt は fmt.Sprintf の簡略ラッパです。
// フォーマット文字列をより簡潔に記述したい場合に使用します。
//
// 例:
//
//	util.Fmt("ui.profile.age_group.%d", 20) → "ui.profile.age_group.20"
func Fmt(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// NullableString は、空文字列の場合 nil を返し、
// 非空文字列の場合はポインタを返すユーティリティです。
func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// DerefString は、nil ポインタの場合は "" を返し、
// 値があればその内容を返す逆変換です。
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
