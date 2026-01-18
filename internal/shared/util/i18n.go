// Package util は i18n 関連の共通ヘルパー関数を提供します。
// 旧 shared/i18n/loader.go から移植された wrapWithPath, deepMerge, flatten などを含みます。
package util

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadYAMLFile はファイルパスを受け取り、YAML を map[string]any にパースします。
func ReadYAMLFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	out := map[string]any{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}
	return out, nil
}

// WrapWithPath は ["ui","page","home"] + {"title":"Home"} を
// {"ui":{"page":{"home":{"title":"Home"}}}} の形に変換します。
func WrapWithPath(prefix []string, in map[string]any) map[string]any {
	if len(prefix) == 0 {
		return in
	}
	last := prefix[len(prefix)-1]
	if len(prefix) == 1 {
		return map[string]any{last: in}
	}
	return map[string]any{prefix[0]: WrapWithPath(prefix[1:], in)}
}

// DeepMerge は src の内容を dst に再帰的にマージします。
// ネスト構造が衝突する場合、map[string]any 同士なら再帰的に統合します。
func DeepMerge(dst, src map[string]any) error {
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			// 両方 map[string]any の場合のみ再帰マージ
			dm, dok := dv.(map[string]any)
			sm, sok := v.(map[string]any)
			if dok && sok {
				if err := DeepMerge(dm, sm); err != nil {
					return err
				}
				continue
			}
		}
		dst[k] = v
	}
	return nil
}

// Flatten はネスト構造を "a.b.c" 形式のキーでフラット化します。
func Flatten(prefix string, in map[string]any, out map[string]any) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch vv := v.(type) {
		case map[string]any:
			Flatten(key, vv, out)
		case map[any]any: // yaml.v2系がこれを返すことがある
			m2 := make(map[string]any, len(vv))
			for kk, vv2 := range vv {
				if s, ok := kk.(string); ok {
					m2[s] = vv2
				}
			}
			Flatten(key, m2, out)
		default:
			out[key] = vv
		}
	}
}

// DetectLocale はパスからロケールコード (ja, en, etc.) を抽出します。
func DetectLocale(path string) string {
	normalized := strings.ReplaceAll(path, `\`, "/")
	segs := strings.Split(normalized, "/")
	for _, cand := range []string{"ja", "en"} {
		for _, s := range segs {
			if s == cand {
				return cand
			}
		}
	}
	return "ja"
}
