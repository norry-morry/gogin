// Package profile は、ユーザープロフィールに関するユースケースロジックを提供します。
//
// このパッケージはアプリケーション層に属し、
// エンティティやリポジトリを利用してドメイン操作を実現します。
// 特に、ユーザープロフィールの参照・更新（PATCH）・登録（Upsert）などの
// ビジネスロジックを実装します。
package profile

import (
	"context"
	"fmt"
	"time"

	"resume/internal/domain/entity"
)

// PatchUserProfile は、ユーザープロフィール情報の登録または更新（Upsert）を行うユースケースです。
//
// 呼び出し元（controller）は、ユーザーIDと部分的なプロフィール情報を
// PatchUserProfileInput 構造体として渡します。
//
// 本メソッドの主な責務は以下の通りです：
// 1. 入力値の検証と正規化（例：BirthDate の日付パース）
// 2. 既存レコードの取得（存在しない場合は新規作成として扱う）
// 3. GenderID（性別ID）の決定：
// - 入力で指定されていればそれを使用
// - 既存レコードがあれば既存値を維持
// - どちらも無ければ「未回答(4)」をセット
// 4. エンティティを構築して Upsert（INSERT ... ON DUPLICATE KEY UPDATE）を実行
//
// 成功時は nil を返し、更新結果は返却しません。
// （参照API側で最新状態を取得する設計としています）
func (uc *Interactor) PatchUserProfile(ctx context.Context, in PatchUserProfileInput) error {
	if in.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	// Entity 生成
	up := &entity.UserProfile{
		UserID:         in.UserID,
		FamilyName:     in.FamilyName,
		GivenName:      in.GivenName,
		FamilyNameKana: in.FamilyNameKana,
		GivenNameKana:  in.GivenNameKana,
	}

	// ==========
	// BirthDate（任意）
	// ==========
	if in.BirthDate.IsCleared() {
		// null 指定 → 「NULL にクリア」したいのでゼロ値ポインタを立てる
		var z time.Time
		up.BirthDate = &z
	} else if in.BirthDate.IsSet() {
		v := *in.BirthDate.Value
		if v == "" {
			// 空文字もクリア扱いにするなら同様
			var z time.Time
			up.BirthDate = &z
		} else {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				return fmt.Errorf("invalid birth_date format: must be YYYY-MM-DD")
			}
			up.BirthDate = &t
		}
	}
	// IsUnset のときは up.BirthDate は nil のまま（＝触らない）

	// ==========
	// GenderID（任意・NULL 可）
	// ==========
	if in.GenderID.IsCleared() {
		// null → クリア（NULL）を表すため 0 をセンチネルに
		var z uint8 = 0
		up.GenderID = &z
	} else if in.GenderID.IsSet() {
		up.GenderID = in.GenderID.Value // 1,2,... の実値
	}
	// IsUnset のときは up.GenderID は nil のまま

	// ==========
	// Initial（任意）
	// ==========
	if in.Initial.IsCleared() {
		// null → クリア（NULL）を表すため空文字をセンチネルに
		empty := ""
		up.Initial = &empty
	} else if in.Initial.IsSet() {
		v := *in.Initial.Value
		if v == "" {
			// 空文字はクリア扱いで NULL にしたいなら同様
			empty := ""
			up.Initial = &empty
		} else {
			up.Initial = &v
		}
	}

	// Upsert 実行
	if err := uc.upRepo.Upsert(ctx, up); err != nil {
		return fmt.Errorf("failed to upsert user_profile: %w", err)
	}

	return nil
}
