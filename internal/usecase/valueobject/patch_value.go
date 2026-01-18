// Package valueobject は、ユースケース層で使用する補助的な値オブジェクトを提供します。
package valueobject

// 特に PATCH リクエストなどにおいて、フィールドの状態を3値で表現するための
// PatchValue 型を定義します。
//
// PatchValue は以下の3つの状態を区別できます：
//
// - JSON にフィールドが **含まれていない** → Unset（更新対象外）
// - JSON にフィールドが **null として含まれる** → Cleared（クリア／リセット）
// - JSON にフィールドが **具体的な値で含まれる** → Set（更新対象）
//
// これにより、「未送信」「明示的に削除」「値を変更」の区別を
// バックエンド側で安全に扱うことができます。

import "encoding/json"

// PatchValue は、PATCH リクエストなどで利用する「3値入力」を表すジェネリックな値オブジェクトです。
// 型パラメータ T は、json.Unmarshal 可能な任意の型を指定できます。
//
// 例えば、次のようにユースケース入力構造体で使用します：
//
//	type PatchUserProfileInput struct {
//		 Initial   valueobject.PatchValue[string] `json:"initial"`
//		 BirthDate valueobject.PatchValue[string] `json:"birthDate"`
//		 GenderID  valueobject.PatchValue[uint8]  `json:"genderId"`
//	}
//
// 各フィールドは次のように評価されます：
//
//   - JSON にキーが存在しない場合 … Set=false, Value=nil （更新しない）
//   - JSON に null が指定された場合 … Set=true, Value=nil （値をクリアする）
//   - JSON に具体的な値が指定された場合 … Set=true, Value=&値 （値を更新する）
//
// これにより、Interactor 側で「未指定・クリア・更新」を安全に判定できます。
type PatchValue[T any] struct {
	Set   bool // JSON にフィールドが存在する場合 true
	Value *T   // 値（null の場合は nil）
}

// UnmarshalJSON は json.Unmarshaler インタフェースを実装します。
// JSON 解析時に「キーが存在したかどうか(Set)」と「値が null かどうか」を判定します。
func (p *PatchValue[T]) UnmarshalJSON(data []byte) error {
	p.Set = true

	// "null" の場合は明示的にクリア指定
	if string(data) == "null" {
		p.Value = nil
		return nil
	}

	// それ以外は通常の値としてパース
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	p.Value = &v
	return nil
}

// IsUnset は、このフィールドが JSON に含まれていなかった場合に true を返します。
// この場合は更新対象外として扱うのが一般的です。
func (p PatchValue[T]) IsUnset() bool {
	return !p.Set
}

// IsCleared は、このフィールドが JSON に含まれ、かつ値が null であった場合に true を返します。
// この場合は、既存の値をクリア（NULL や空文字に更新）することを意図します。
func (p PatchValue[T]) IsCleared() bool {
	return p.Set && p.Value == nil
}

// IsSet は、このフィールドが JSON に含まれ、かつ値が null ではない場合に true を返します。
// この場合は、新しい値に更新することを意図します。
func (p PatchValue[T]) IsSet() bool {
	return p.Set && p.Value != nil
}
