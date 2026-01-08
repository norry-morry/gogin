// Package db は、GORM を利用したデータベース接続および
// トランザクション制御の共通ユーティリティを提供します。
//
// 本ファイル (tx.go) は、usecase 層で定義される TxRunner インターフェースの
// GORM 実装を提供し、トランザクションを context.Context 経由で伝播させます。
//
// --- 概要 Summary ---
//   - GormTxRunner.Do() : トランザクションを開始して ctx に注入し、fn を実行
//   - FromCtxOrDB()           : ctx に含まれる tx（*gorm.DB）を取り出す。無ければ通常DB
//
// この設計により、usecase 層で Transaction 境界を一元管理し、
// repository 層では「ctxからtxを取り出して使う」だけで済むようになります。
//
// これにより、users と auth_identities など複数テーブル更新が1つのトランザクションで安全に実行されます。
package db

import (
	"context"

	"gorm.io/gorm"
)

// txKey は context.Context に格納するトランザクション識別キーです。
// context の Code 検索時に衝突しないよう、独自の非公開型を使用します。
type txKey struct {
}

// GormTxRunner は GORM 用のトランザクション実行構造体です。
// usecase 層の TxRunner インターフェースを実装します。
type GormTxRunner struct {
	db *gorm.DB
}

// NewTxRunner は GormTxRunner のコンストラクタです。
// DI (wire) などから呼び出して usecase 層に注入します。
func NewTxRunner(db *gorm.DB) *GormTxRunner {
	return &GormTxRunner{
		db: db,
	}
}

// Do は与えられた関数 fn をトランザクション内で実行します。
// ctx に tx(*gorm.DB) を格納し、repository 層からは同一トランザクションを参照できます。
// エラーが返された場合はロールバックし、nil の場合はコミットします。
func (t *GormTxRunner) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx2 := context.WithValue(ctx, txKey{}, tx)
		return fn(ctx2)
	})
}

// FromCtxOrDB は、context からトランザクション (*gorm.DB) を取得します。
// ctx にトランザクションが含まれていればそれを返し、無ければ通常の db 接続を返します。
// repository 層ではこの関数を利用することで、
// トランザクション境界を意識せずに安全に db を扱えます。
func FromCtxOrDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return db.WithContext(ctx)
}

// WithTx は お好みで：Txヘルパ
func WithTx(ctx context.Context, db *gorm.DB, fn func(ctx context.Context) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(InjectTx(ctx, tx))
	})
}

// InjectTx は gorm.DB のトランザクションを context.Context に埋め込みます。
// 後続のリポジトリ層などで、context から現在のトランザクションを取得できるようにします。
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
