// Package errors は go-playground/validator/v10 の ValidationErrors を
// API レスポンスの details 形式（フィールドごとのタグ/パラメータ/値）へ変換する
// ユーティリティを提供します。
package errors

import (
	"errors"

	v10 "github.com/go-playground/validator/v10"
)

// ToDetails は、バリデーションエラーを API レスポンスの details 形式に変換します。
// 各フィールドに対応するタグ名（tag）、パラメータ（param）、入力値（value）を含むマップを返します。
// ValidationErrors でない場合は、エラーメッセージを単一要素のマップ {"error": "..."} として返します。
func ToDetails(err error) map[string]any {
	var verr v10.ValidationErrors
	if !errors.As(err, &verr) {
		return map[string]any{"error": err.Error()}
	}
	out := make(map[string]any, len(verr))
	for _, fe := range verr {
		// fe.Field() は構造体のフィールド名（JSON名にしたければタグから引く実装を足す）
		out[fe.Field()] = map[string]any{
			"tag":   fe.Tag(),   // e.g. "required", "max", "jp_postal"
			"param": fe.Param(), // e.g. "160"
			"value": fe.Value(),
		}
	}
	return out
}
