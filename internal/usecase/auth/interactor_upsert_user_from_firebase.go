// Package auth は、認証ユースケース（Usecase）の実装を提供します。
// 本ファイルでは「Firebase IDトークン検証済みのユーザー情報を元に、
// users / auth_identities を upsert し、レスポンスDTOを返す」ユースケースを実装します。
package auth

import (
	"context"
	"fmt"

	"resume/internal/shared/apperr"
	"resume/internal/shared/util"
)

// UpsertUserFromFirebase は、VerifyIDToken 済みの入力（UID / プロバイダ配列 等）を受け取り、
// - users テーブルの作成 or 更新（LastLoginAt をサーバ時刻で更新）
// - auth_identities テーブルの upsert（provider, provider_user_id で一意）
// を **1トランザクション** で実行します。
// 成功時は、最新のユーザー情報（紐付くプロバイダ一覧を含む）とトークンを返します。
//
// エラーポリシー：
// - 入力不足（IDToken/UID 無し）は ErrInvalidInput を返す
// - リポジトリ層のDBエラーは mapInfraToUCError で UC語彙（ErrConflict/ErrNotFound 等）に正規化しつつ伝播
// - トランザクション内のいずれかで失敗したらその場で中断・ロールバック
func (u *interactor) UpsertUserFromFirebase(ctx context.Context, in UpsertUserFromFirebaseInput) (AuthResultOutput, error) {
	// --- 0) 入力チェック -------------------------------------------------------
	// ここでは必須だけを軽く確認（詳細バリデーションは上位レイヤで実施しても良い）。
	if in.IDToken == "" || in.UID == "" {
		return AuthResultOutput{}, apperr.New(
			apperr.CodeUnprocessable,
			"validation failed",
			map[string]any{"fields": map[string]string{
				"idToken": "required",
				"uid":     "required",
			}},
		)
	}

	// UCに注入された Clock を用いて現在時刻（UTC想定）を取得。
	// これを LastLoginAt / CreatedAt / UpdatedAt に用いる。
	now := u.clock.Now()

	// 以降で使う変数を先に宣言。
	var (
		isNew   bool   // 今回作成された新規ユーザーかどうか
		userEnt *User  // 処理の最終的な対象ユーザー（作成 or 更新後）
		userID  uint64 // auth_identities 側のFKに合わせた int64 のID
	)

	// --- 1) トランザクション境界 ------------------------------------------------
	// TxRunner.Do は ctx に *gorm.db(tx) を差し込み、repo が FromCtxOrDB(...) で
	// その tx を取得して同一トランザクションに乗れるようにする。
	if err := u.tx.Do(ctx, func(txCtx context.Context) error {
		// --- 2) users の upsert ------------------------------------------------
		// UID は Firebase プロジェクト内で一意。まず既存のユーザーを探す。
		found, err := u.users.FindByUID(txCtx, in.UID)
		if err != nil {
			// DBエラーはそのまま返し、Tx もロールバック。
			return fmt.Errorf("find user by uid: %w", err)
		}

		if found == nil {
			// 2-a) 見つからなければ新規作成
			isNew = true
			userEnt = &User{
				UID: in.UID,
				// util.Clone は *T をシャローコピー。nil 安全に上書きしたい時に使う。
				Email:         util.Clone(in.Email),
				EmailVerified: in.EmailVerified,
				DisplayName:   util.Clone(in.DisplayName),
				PhotoURL:      util.Clone(in.PhotoURL),
				LastLoginAt:   &now,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			newID, err := u.users.Create(txCtx, userEnt)
			if err != nil {
				// DB一意制約などは repo 側で整形 → ここで UC 語彙に正規化
				return mapInfraToUCError(fmt.Errorf("create user: %w", err))
			}
			// domain/entity.User は ID が uint64 の想定。FK用に int64 も保持。
			userEnt.ID = newID
			userID = newID
		} else {
			// 2-b) 見つかったのでプロフィール同期＋最終ログイン更新
			found.Email = util.Clone(in.Email)
			found.EmailVerified = in.EmailVerified
			found.DisplayName = util.Clone(in.DisplayName)
			found.PhotoURL = util.Clone(in.PhotoURL)
			found.LastLoginAt = &now
			found.UpdatedAt = now

			// プロフィールのログイン時更新はリポジトリ側でカラム限定 Updates。
			if err := u.users.UpdateProfileOnLogin(txCtx, found); err != nil {
				return mapInfraToUCError(fmt.Errorf("update user on login: %w", err))
			}
			userEnt = found
			userID = found.ID
		}

		// --- 3) auth_identities の upsert -------------------------------------
		// Providers は 0..n を許容。password/anonymous なども来うる。
		// provider, provider_user_id の組で一意にし、ON CONFLICT DO UPDATE の想定。
		for _, p := range in.Providers {
			ident := &Identity{
				UserID:              userID,
				Provider:            p.Provider,
				ProviderUserID:      p.ProviderUserID,
				ProviderDisplayName: p.ProviderDisplayName, // NOT NULL 前提。ハンドラでFallback済みを想定。
				EmailAtSignup:       util.Clone(p.EmailAtSignup),
				CreatedAt:           now, // 新規時採用。既存更新時は repo 実装が保持。
			}
			if _, err := u.idents.Upsert(txCtx, ident); err != nil {
				// (provider, provider_user_id) の重複やFK不整合などをUC語彙へマップ
				return mapInfraToUCError(fmt.Errorf("upsert identity (%s/%s): %w", p.Provider, p.ProviderUserID, err))
			}
		}
		return nil
	}); err != nil {
		// Tx 内で発生したエラーはここに伝播。HTTPハンドラは標準のエラーフォーマットで返す。
		return AuthResultOutput{}, err
	}

	// --- 4) 返却用のプロバイダ一覧を再読込 --------------------------------------
	// Tx外で読み出してよい（副作用なし）。一覧を DTO に詰め替えて返す。
	idents, err := u.idents.ListByUserID(ctx, userID)
	if err != nil {
		return AuthResultOutput{}, mapInfraToUCError(fmt.Errorf("list providers: %w", err))
	}

	// --- 5) 出力DTOを返却 ------------------------------------------------------
	// Token は今回は受け取った IDToken をそのまま返す想定。
	return AuthResultOutput{
		Token:     in.IDToken,
		User:      toUserOutput(userEnt, idents),
		IsNewUser: isNew,
	}, nil
}
