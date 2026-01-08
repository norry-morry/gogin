// Package jpstring は、日本語文字列を正規化するユーティリティ関数群を提供します。
// 半角／全角の揺れや空白・記号の統一など、バリデーション前の前処理を目的とします。
package jpstring

import (
	"strings"
	"unicode"
)

var (
	fullToHalfNum = map[rune]rune{
		'０': '0', '１': '1', '２': '2', '３': '3', '４': '4', '５': '5', '６': '6', '７': '7', '８': '8', '９': '9',
	}
	hyphenVariants = []rune{'ー', '－', '―', '‐', '−'}
)

// NormalizeSpaces は、全角スペースを半角スペースに変換し、
// 連続する空白を 1 つにまとめて返します。
func NormalizeSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '　'
	})
}

// NormalizeDigitsHyphen は、全角数字と全角ハイフンを半角に統一します。
func NormalizeDigitsHyphen(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if h, ok := fullToHalfNum[r]; ok {
			out = append(out, h)
			continue
		}
		replaced := false
		for _, hv := range hyphenVariants {
			if r == hv {
				out = append(out, '-')
				replaced = true
				break
			}
		}
		if replaced {
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

// NormalizeJP は、日本語文字列全体を正規化します。
// 全角／半角、スペース、数字・ハイフンの統一などを行います。
func NormalizeJP(s string) string {
	return NormalizeDigitsHyphen(NormalizeSpaces(s))
}
