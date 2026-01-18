// Package valerr は、ドメイン層で発生する値検証エラーを表現するための共通インターフェースおよび型定義を提供します。
// バリデーションエラーの集約やエラーメッセージ整形などに利用します。
package valerr

import "resume/internal/domain/entity"

// FromDomain converts domain-level InvalidAddressError to a handler/presenter friendly details map.
// Same shape as your handlerutil's "details".
func FromDomain(inv entity.InvalidAddressError) map[string]any {
	if len(inv.Problems) == 0 {
		return nil
	}
	out := make(map[string]any, len(inv.Problems))
	for _, p := range inv.Problems {
		// last-wins（同一フィールドに複数ある場合は最後を採用）
		out[p.Field] = map[string]any{
			"tag":   p.Tag,
			"param": p.Param,
		}
	}
	return out
}

// FromDomainAll converts domain-level InvalidAddressError to a handler/presenter friendly details map.
// Same shape as your handlerutil's "details".
// “同一フィールドに複数エラーを全部出す”版
func FromDomainAll(inv entity.InvalidAddressError) map[string]any {
	if len(inv.Problems) == 0 {
		return nil
	}
	out := make(map[string]any)
	for _, p := range inv.Problems {
		d := map[string]any{"tag": p.Tag, "param": p.Param}
		if cur, ok := out[p.Field]; ok {
			switch vv := cur.(type) {
			case []map[string]any:
				out[p.Field] = append(vv, d)
			case map[string]any:
				out[p.Field] = []map[string]any{vv, d}
			}
		} else {
			out[p.Field] = d
		}
	}
	return out
}
