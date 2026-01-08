// Package dberr_test は MySQL エラー変換ユーティリティ
// (internal/infra/db/dberr/mysql.go) のユニットテストを提供します。
package dberr_test

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"resume/internal/infra/db/dberr"
	"resume/internal/shared/apperr"
)

//
// ---------- nil のとき ----------
//

func TestMap_Nil(t *testing.T) {
	t.Parallel()

	if got := dberr.Map(nil); got != nil {
		t.Errorf("expected nil, got=%v", got)
	}
}

//
// ---------- GORM エラー系 ----------
//

func TestMap_GormErrDuplicatedKey(t *testing.T) {
	t.Parallel()

	err := gorm.ErrDuplicatedKey
	got := dberr.Map(err)

	appErr, ok := got.(*apperr.Error)
	if !ok {
		t.Fatalf("expected *apperr.Error, got=%T", got)
	}
	if appErr.Code != apperr.CodeConflict {
		t.Errorf("expected CodeConflict, got=%s", appErr.Code)
	}
}

func TestMap_GormErrNotFound(t *testing.T) {
	t.Parallel()

	err := gorm.ErrRecordNotFound
	got := dberr.Map(err)

	appErr, ok := got.(*apperr.Error)
	if !ok {
		t.Fatalf("expected *apperr.Error, got=%T", got)
	}
	if appErr.Code != apperr.CodeNotFound {
		t.Errorf("expected CodeNotFound, got=%s", appErr.Code)
	}
}

//
// ---------- MySQL Error Number 系 ----------
//

func newMysqlErr(code uint16) error {
	return &mysql.MySQLError{
		Number: code,
	}
}

func expectCode(t *testing.T, err error, expected apperr.Code) {
	t.Helper()

	appErr, ok := err.(*apperr.Error)
	if !ok {
		t.Fatalf("expected *apperr.Error, got=%T", err)
	}
	if appErr.Code != expected {
		t.Errorf("expected code=%s, got=%s", expected, appErr.Code)
	}
}

// 1062: Duplicate entry
func TestMap_Mysql_DuplicateEntry(t *testing.T) {
	t.Parallel()

	got := dberr.Map(newMysqlErr(1062))
	expectCode(t, got, apperr.CodeConflict)
}

// 1452 / 1216 / 1217: Foreign key constraint
func TestMap_Mysql_ForeignKey1(t *testing.T) {
	t.Parallel()
	got := dberr.Map(newMysqlErr(1452))
	expectCode(t, got, apperr.CodeUnprocessable)
}

func TestMap_Mysql_ForeignKey2(t *testing.T) {
	t.Parallel()
	got := dberr.Map(newMysqlErr(1216))
	expectCode(t, got, apperr.CodeUnprocessable)
}

func TestMap_Mysql_ForeignKey3(t *testing.T) {
	t.Parallel()
	got := dberr.Map(newMysqlErr(1217))
	expectCode(t, got, apperr.CodeUnprocessable)
}

// 1048: NOT NULL violation
func TestMap_Mysql_NullConstraint(t *testing.T) {
	t.Parallel()

	got := dberr.Map(newMysqlErr(1048))
	expectCode(t, got, apperr.CodeUnprocessable)
}

// 3819: CHECK constraint violation
func TestMap_Mysql_CheckConstraint(t *testing.T) {
	t.Parallel()

	got := dberr.Map(newMysqlErr(3819))
	expectCode(t, got, apperr.CodeUnprocessable)
}

//
// ---------- fallback ----------
//

func TestMap_Fallback(t *testing.T) {
	t.Parallel()

	original := errors.New("some other error")
	got := dberr.Map(original)

	// 変換されず同じエラーが返ること
	if !errors.Is(got, original) {
		t.Errorf("expected original error to be returned, got=%v", got)
	}
}
