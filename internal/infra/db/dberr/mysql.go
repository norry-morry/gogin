// Package dberr は、MySQL のエラーをアプリケーション共通のエラーコードへ変換します。
package dberr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"resume/internal/shared/apperr"
)

// Map は GORM/MySQL のエラーをアプリ共通のエラーに変換する
func Map(err error) error {
	if err == nil {
		return nil
	}

	// まずGORMの汎用エラー
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperr.New(apperr.CodeConflict, "duplicate resource", nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.New(apperr.CodeNotFound, "record not found", nil)
	}

	// MySQLドライバの番号で詳細に振り分け
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		case 1062: // ER_DUP_ENTRY
			return apperr.New(apperr.CodeConflict, "duplicate resource", nil)
		case 1452, 1216, 1217: // FK fails / cannot add/update/delete due to FK
			return apperr.New(apperr.CodeUnprocessable, "invalid foreign key", nil)
		case 1048: // ER_BAD_NULL_ERROR (NOT NULL 列に NULL)
			return apperr.New(apperr.CodeUnprocessable, "null constraint violation", nil)
		case 3819: // CHECK constraint
			return apperr.New(apperr.CodeUnprocessable, "check constraint violation", nil)
		}
	}
	// 既定: そのまま返す（上位で500へ）
	return err
}
