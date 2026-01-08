// Package util は、共通の軽量ユーティリティ関数群を提供します。
package util

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// 機微情報マスキングユーティリティ
//
// このパッケージは、リクエストやレスポンスなどの構造化ログ内に含まれる
// パスワード・トークン・個人情報等を安全にマスクするための共通ユーティリティを提供する。
// YAML 設定ファイルからマスク対象キーを読み込み、キー名に応じて完全マスク・部分マスクを自動適用する。

// SensitiveConfig は YAML から読み込まれるマスク設定情報を表す構造体。
type SensitiveConfig struct {
	Sensitive struct {
		Full []string `yaml:"full"` // 完全マスク対象キーの一覧
		Part []string `yaml:"part"` // 部分マスク対象キーの一覧
	} `yaml:"sensitive"`
	Partial struct {
		KeepHead int    `yaml:"keep_head"` // 部分マスク時に残す先頭文字数
		KeepTail int    `yaml:"keep_tail"` // 部分マスク時に残す末尾文字数
		MaskChar string `yaml:"mask_char"` // マスクに使用する文字
	} `yaml:"partial"`
}

// Masker は機微情報マスカーの抽象インタフェース。
// 具体実装（keyMatcher）は外部に公開しない。
type Masker interface {
	MaskByKey(key string, val any) any
	MaskMapShallow(src map[string]any) map[string]any
	MaskAnyRecursive(key string, v any) any
}

// NewMasker は設定から Masker を生成するファサード。
// 具体型 keyMatcher を隠蔽したい場合はこの関数を使う。
func NewMasker(cfg *SensitiveConfig) (Masker, error) {
	return NewKeyMatcher(cfg)
}

// LoadSensitiveConfig は YAML ファイルからマスク設定を読み込む。
// ファイルが存在しない場合やパースに失敗した場合はエラーを返す。
func LoadSensitiveConfig(path string) (*SensitiveConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c SensitiveConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	// デフォルト値を補完
	if c.Partial.KeepHead < 0 {
		c.Partial.KeepHead = 1
	}
	if c.Partial.KeepTail < 0 {
		c.Partial.KeepTail = 4
	}
	if c.Partial.MaskChar == "" {
		c.Partial.MaskChar = "*"
	}

	return &c, nil
}

// Redacted は完全マスク時に使用される固定文字列。
const Redacted = "[REDACTED]"

// FullMask は渡された値を完全にマスクし、固定文字列 [REDACTED] を返す。
func FullMask(_ any) string {
	return Redacted
}

// PartialMask は渡された文字列を部分的にマスクする。
// 先頭と末尾を残し、中間を指定のマスク文字で埋める。
func PartialMask(s string, keepHead, keepTail int, maskChar string) string {
	if s == "" {
		return s
	}
	if keepHead < 0 {
		keepHead = 0
	}
	if keepTail < 0 {
		keepTail = 0
	}
	if maskChar == "" {
		maskChar = "*"
	}

	runes := []rune(s)
	n := len(runes)
	if keepHead+keepTail >= n {
		// 残す長さが総文字数以上なら、そのまま返す
		return s
	}

	head := string(runes[:keepHead])
	tail := string(runes[n-keepTail:])
	midLen := n - keepHead - keepTail

	return head + strings.Repeat(maskChar, midLen) + tail
}

// keyMatcher はキー名に対するマスク種別（完全／部分）判定を行う構造体。
type keyMatcher struct {
	full []*regexp.Regexp // 完全マスク対象の正規表現リスト
	part []*regexp.Regexp // 部分マスク対象の正規表現リスト

	keepHead int
	keepTail int
	maskChar string
}

// globToRegex はワイルドカード（*）を含むキー定義を正規表現に変換する。
// 大文字小文字は区別しない。
func globToRegex(glob string) (*regexp.Regexp, error) {
	p := regexp.QuoteMeta(glob)
	p = strings.ReplaceAll(p, "\\*", ".*")
	p = "(?i)^.*" + p + ".*$" // キー名部分一致許可
	return regexp.Compile(p)
}

// NewKeyMatcher は SensitiveConfig から keyMatcher を生成する。
// Full/Part のパターンを正規表現化して保持する。
func NewKeyMatcher(cfg *SensitiveConfig) (*keyMatcher, error) {
	if cfg == nil {
		return nil, errors.New("nil config")
	}

	m := &keyMatcher{
		keepHead: cfg.Partial.KeepHead,
		keepTail: cfg.Partial.KeepTail,
		maskChar: cfg.Partial.MaskChar,
	}

	// 完全マスクパターンをコンパイル
	for _, g := range cfg.Sensitive.Full {
		re, err := globToRegex(g)
		if err != nil {
			return nil, fmt.Errorf("compile full pattern %q: %w", g, err)
		}
		m.full = append(m.full, re)
	}

	// 部分マスクパターンをコンパイル
	for _, g := range cfg.Sensitive.Part {
		re, err := globToRegex(g)
		if err != nil {
			return nil, fmt.Errorf("compile part pattern %q: %w", g, err)
		}
		m.part = append(m.part, re)
	}

	return m, nil
}

// isFullKey は与えられたキーが完全マスク対象に一致するかを判定する。
func (m *keyMatcher) isFullKey(key string) bool {
	for _, re := range m.full {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// isPartKey は与えられたキーが部分マスク対象に一致するかを判定する。
func (m *keyMatcher) isPartKey(key string) bool {
	for _, re := range m.part {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskByKey はキー名に応じて値を完全または部分的にマスクする。
// 対象外のキーはそのまま返す。
func (m *keyMatcher) MaskByKey(key string, val any) any {
	if key == "" {
		return val
	}
	if m.isFullKey(key) {
		return FullMask(val)
	}
	if m.isPartKey(key) {
		s := fmt.Sprintf("%v", val)
		return PartialMask(s, m.keepHead, m.keepTail, m.maskChar)
	}
	return val
}

// MaskMapShallow は map[string]any を対象に、
// 各キーに応じてマスク処理を適用する（1階層のみ）。
func (m *keyMatcher) MaskMapShallow(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = m.MaskByKey(k, v)
	}
	return dst
}

// MaskAnyRecursive は任意の構造（map やスライスを含む）に対して
// 再帰的にマスク処理を適用する。
// JSON デコード結果などネストした構造体を安全にマスクしたい場合に使用する。
func (m *keyMatcher) MaskAnyRecursive(key string, v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, vv := range t {
			out[k] = m.MaskAnyRecursive(k, vv)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, vv := range t {
			out[i] = m.MaskAnyRecursive(key, vv)
		}
		return out
	default:
		return m.MaskByKey(key, v)
	}
}
