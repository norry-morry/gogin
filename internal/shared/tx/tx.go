// Package tx はユースケース横断で使うトランザクションI/Fを定義します。
package tx

import "context"

// Runner は「ctxにTXを伝搬しつつfnをTX内で実行する」共通I/F
// メソッド名は Do に統一
type Runner interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
