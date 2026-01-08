// Package db_test は、GORM を利用したトランザクションユーティリティ
// (internal/infra/db/tx.go) のユニットテストを提供します。
package db_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	dbinfra "resume/internal/infra/db"
)

//
// ---------- テスト用ヘルパ ----------
//

// txUser はトランザクション挙動を検証するためのシンプルなモデルです。
type txUser struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// ★ テスト名ごとに異なる in-memory DB を使う
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	if err := db.AutoMigrate(&txUser{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

//
// ---------- FromCtxOrDB / InjectTx ----------
//

// TestFromCtxOrDB_UsesTxFromContext は、context に埋め込まれた tx が優先されることを確認します。
func TestFromCtxOrDB_UsesTxFromContext(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	ctx := context.Background()

	// 通常 DB とは別の tx を作って context に埋め込む
	err := baseDB.Transaction(func(tx *gorm.DB) error {
		ctxWithTx := dbinfra.InjectTx(ctx, tx)

		got := dbinfra.FromCtxOrDB(ctxWithTx, baseDB)
		if got != tx {
			t.Errorf("expected FromCtxOrDB to return tx from context, got different *gorm.DB")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error in transaction: %v", err)
	}
}

// TestFromCtxOrDB_FallsBackToBaseDB は、context に tx が無い場合に通常の DB を返すことを確認します。
func TestFromCtxOrDB_FallsBackToBaseDB(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	ctx := context.Background()

	// FromCtxOrDB で取得した DB を使って insert → baseDB からも見えることを確認
	dbForUse := dbinfra.FromCtxOrDB(ctx, baseDB)
	if err := dbForUse.Create(&txUser{Name: "fallback"}).Error; err != nil {
		t.Fatalf("failed to create via FromCtxOrDB: %v", err)
	}

	var count int64
	if err := baseDB.Model(&txUser{}).Where("name = ?", "fallback").Count(&count).Error; err != nil {
		t.Fatalf("failed to count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

//
// ---------- GormTxRunner.Do ----------
//

// TestGormTxRunner_Do_CommitsOnNilError は、fn が nil を返した場合にコミットされることを確認します。
func TestGormTxRunner_Do_CommitsOnNilError(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	runner := dbinfra.NewTxRunner(baseDB)

	ctx := context.Background()
	err := runner.Do(ctx, func(ctx context.Context) error {
		// ctx から tx を取り出して insert
		txDB := dbinfra.FromCtxOrDB(ctx, baseDB)
		if err := txDB.Create(&txUser{Name: "committed"}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error from Do, got %v", err)
	}

	var count int64
	if err := baseDB.Model(&txUser{}).Where("name = ?", "committed").Count(&count).Error; err != nil {
		t.Fatalf("failed to count after commit: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 committed row, got %d", count)
	}
}

// TestGormTxRunner_Do_RollbacksOnError は、fn がエラーを返した場合にロールバックされることを確認します。
func TestGormTxRunner_Do_RollbacksOnError(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	runner := dbinfra.NewTxRunner(baseDB)

	ctx := context.Background()
	expectedErr := errors.New("force rollback")

	err := runner.Do(ctx, func(ctx context.Context) error {
		txDB := dbinfra.FromCtxOrDB(ctx, baseDB)
		if err := txDB.Create(&txUser{Name: "rollback"}).Error; err != nil {
			return err
		}
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	var count int64
	if err := baseDB.Model(&txUser{}).Where("name = ?", "rollback").Count(&count).Error; err != nil {
		t.Fatalf("failed to count after rollback: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows after rollback, got %d", count)
	}
}

//
// ---------- WithTx ----------
//

// TestWithTx_CommitsOnNilError は、WithTx が fn=nil エラーでコミットされることを確認します。
func TestWithTx_CommitsOnNilError(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	ctx := context.Background()

	err := dbinfra.WithTx(ctx, baseDB, func(ctx context.Context) error {
		txDB := dbinfra.FromCtxOrDB(ctx, baseDB)
		if err := txDB.Create(&txUser{Name: "withtx-commit"}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error from WithTx, got %v", err)
	}

	var count int64
	if err := baseDB.Model(&txUser{}).Where("name = ?", "withtx-commit").Count(&count).Error; err != nil {
		t.Fatalf("failed to count after commit: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 committed row, got %d", count)
	}
}

// TestWithTx_RollbacksOnError は、WithTx がエラー時にロールバックされることを確認します。
func TestWithTx_RollbacksOnError(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	ctx := context.Background()
	expectedErr := errors.New("withtx rollback")

	err := dbinfra.WithTx(ctx, baseDB, func(ctx context.Context) error {
		txDB := dbinfra.FromCtxOrDB(ctx, baseDB)
		if err := txDB.Create(&txUser{Name: "withtx-rollback"}).Error; err != nil {
			return err
		}
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	var count int64
	if err := baseDB.Model(&txUser{}).Where("name = ?", "withtx-rollback").Count(&count).Error; err != nil {
		t.Fatalf("failed to count after rollback: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows after rollback, got %d", count)
	}
}

//
// ---------- InjectTx 単体 ----------
//

// TestInjectTx_AndFromCtxOrDB は InjectTx で埋め込んだ tx を FromCtxOrDB で取得できることを確認します。
func TestInjectTx_AndFromCtxOrDB(t *testing.T) {
	t.Parallel()

	baseDB := newTestDB(t)
	ctx := context.Background()

	err := baseDB.Transaction(func(tx *gorm.DB) error {
		ctxWithTx := dbinfra.InjectTx(ctx, tx)
		got := dbinfra.FromCtxOrDB(ctxWithTx, baseDB)
		if got != tx {
			t.Errorf("expected FromCtxOrDB to return injected tx, got different *gorm.DB")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error in transaction: %v", err)
	}
}
