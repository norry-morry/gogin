// Package valerr は、ドメイン層で発生する値検証エラーを表現するための共通インターフェースおよび型定義を提供します。
// バリデーションエラーの集約やエラーメッセージ整形などに利用します。
package valerr

// FieldError は 1 フィールドに対する単一の検証エラーを表します。
// JSON では {"code": "...", "params": {...}} の形でシリアライズされます。
type FieldError struct {
	// Code はエラーコードです (例: "validation.rules.required", "validation.rules.duplicate" など)。
	Code string `json:"code"`

	// Params はメッセージ組み立てや i18n に利用する追加パラメータです。
	// 例: {"field": "domain.address.purposeId.label"}
	Params map[string]any `json:"params,omitempty"`
}

// ValidationErrors はフィールド名ごとの検証エラー一覧を表します。
// JSON では {"domain.address.purposeId": [ {..}, {..} ]} のようにシリアライズされます。
type ValidationErrors map[string][]FieldError

// New は空の ValidationErrors を返します。
func New() ValidationErrors {
	return make(ValidationErrors)
}

// Add は指定フィールドに FieldError を追加します。
func (ve ValidationErrors) Add(field, code string, params map[string]any) {
	ve[field] = append(ve[field], FieldError{
		Code:   code,
		Params: params,
	})
}

// ToDetails は app エラーの Details フィールドに設定可能な形式へ変換します。
// apperr.New の第3引数 (details map[string]any) にそのまま渡すために利用します。
func (ve ValidationErrors) ToDetails() map[string]any {
	if ve == nil {
		return nil
	}

	details := make(map[string]any, len(ve))
	for k, v := range ve {
		details[k] = v
	}
	return details
}
