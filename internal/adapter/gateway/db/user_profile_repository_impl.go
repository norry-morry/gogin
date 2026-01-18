// Package db は、GORM を用いたリポジトリ実装を提供するパッケージです。
package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	db2 "resume/internal/infra/db"
)

// userProfileRepository は UserProfileRepository の GORM 実装です。
type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository は UserProfileRepository の新しいインスタンスを生成して返します。
func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

// FindByUserID は、指定されたユーザーIDに紐づくユーザープロフィールを1件取得します。
// 対応するレコードが存在しない場合は、(*entity.UserProfile)(nil), nil を返します。
func (r *userProfileRepository) FindByUserID(ctx context.Context, userID uint64) (*entity.UserProfile, error) {
	db := db2.
		FromCtxOrDB(ctx, r.db).
		Model(&entity.UserProfile{}).
		Preload("Gender")

	var up entity.UserProfile
	err := db.
		Where("user_id = ?", userID).
		Limit(1).
		Take(&up).
		Error

	switch {
	case err == nil:
		return &up, nil

	case errors.Is(err, gorm.ErrRecordNotFound):
		// 対応するレコードが存在しない場合は nil を返す（404 ではなく 200 空レス想定）
		return nil, nil

	default:
		// それ以外のDBエラーはそのまま返す
		return nil, err
	}
}

func (r *userProfileRepository) Upsert(ctx context.Context, p *entity.UserProfile) error {
	if p == nil {
		return fmt.Errorf("nil profile")
	}

	// 必須カラム
	columns := []string{
		"user_id",
		"family_name",
		"given_name",
		"family_name_kana",
		"given_name_kana",
	}
	placeholders := []string{"?", "?", "?", "?", "?"}
	args := []interface{}{
		p.UserID,
		p.FamilyName,
		p.GivenName,
		p.FamilyNameKana,
		p.GivenNameKana,
	}

	// 任意カラム: nil でなければ INSERT/UPDATE 対象に含める
	// 任意カラム: ポインタが nil なら触らない
	// ポインタが立っていたら、「ゼロ値なら NULL」「それ以外なら値そのもの」
	if p.BirthDate != nil {
		columns = append(columns, "birth_date")
		placeholders = append(placeholders, "?")
		if p.BirthDate.IsZero() {
			args = append(args, nil) // NULL で更新
		} else {
			args = append(args, p.BirthDate)
		}
	}

	if p.GenderID != nil {
		columns = append(columns, "gender_id")
		placeholders = append(placeholders, "?")
		if *p.GenderID == 0 {
			args = append(args, nil) // NULL で更新（未回答）
		} else {
			args = append(args, p.GenderID)
		}
	}

	if p.Initial != nil {
		columns = append(columns, "initial")
		placeholders = append(placeholders, "?")
		if *p.Initial == "" {
			args = append(args, nil) // NULL で更新
		} else {
			args = append(args, p.Initial)
		}
	}

	// ON DUPLICATE KEY UPDATE 句を組み立てる
	// 「指定されたカラムだけ UPDATE」したいので、INSERT 対象に入れたカラムから user_id だけ除外
	setClauses := make([]string, 0, len(columns))
	for _, col := range columns {
		if col == "user_id" {
			continue // PK は更新しない
		}
		// MySQL: VALUES(col) で INSERT 側の値を使って UPDATE する
		setClauses = append(setClauses, fmt.Sprintf("%s = VALUES(%s)", col, col))
	}

	sql := fmt.Sprintf(
		"INSERT INTO user_profiles (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(setClauses, ", "),
	)

	return db2.
		FromCtxOrDB(ctx, r.db).
		Model(&entity.UserProfile{}).
		Exec(sql, args...).
		Error
}
