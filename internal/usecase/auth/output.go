// Package auth は、認証関連ユースケースで使用する入出力DTOを定義します。
package auth

import "time"

// ProviderOutput は、auth_identities の1レコードに対応する返却用DTOです。
// アプリ側では「連携中のサインインプロバイダ情報」として利用します。
type ProviderOutput struct {
	ID                  uint64    `json:"id"`                        // auth_identities.id
	Provider            string    `json:"provider"`                  // 例: "google.com", "github.com", "password"
	ProviderUserID      string    `json:"provider_user_id"`          // プロバイダ側の一意ID
	ProviderDisplayName string    `json:"provider_display_name"`     // プロバイダでの表示名
	EmailAtSignup       *string   `json:"email_at_signup,omitempty"` // サインアップ時のメール（NULL許容）
	LinkedAt            time.Time `json:"linked_at"`                 // 連携日時（auth_identities.created_at）
	UpdatedAt           time.Time `json:"updated_at"`                // 最終更新日時（auth_identities.updated_at）
}

// UserOutput は、users テーブルに対応する返却用DTOです。
// フロント側がプロフィールや状態を表示・保持するのに必要な項目をまとめます。
type UserOutput struct {
	ID            uint64           `json:"id"`                      // users.id
	UID           string           `json:"uid"`                     // Firebase UID（プロジェクト内一意）
	Email         *string          `json:"email,omitempty"`         // NULL許容
	EmailVerified bool             `json:"email_verified"`          // メール認証済みか
	DisplayName   *string          `json:"display_name,omitempty"`  // 表示名（NULL許容）
	PhotoURL      *string          `json:"photo_url,omitempty"`     // アイコンURL（NULL許容）
	LastLoginAt   *time.Time       `json:"last_login_at,omitempty"` // 直近ログイン（サーバー時刻）NULLの場合あり
	CreatedAt     time.Time        `json:"created_at"`              // 作成日時
	UpdatedAt     time.Time        `json:"updated_at"`              // 更新日時
	Providers     []ProviderOutput `json:"providers"`               // 連携中のプロバイダ一覧
}

// AuthResultOutput は、VerifyIDToken 成功後の upsert 完了時に返却するレスポンスDTOです。
// アクセストークン（今回は ID トークンをそのまま返す想定）と、ユーザー情報を含みます。
// revive:disable:exported  // パッケージ名に "auth" を残したい設計のため、命名警告を抑制
type AuthResultOutput struct {
	Token     string     `json:"token"`       // 返却するトークン（例: Firebase ID トークン）
	User      UserOutput `json:"user"`        // ユーザー情報
	IsNewUser bool       `json:"is_new_user"` // 今回の処理でユーザーを新規作成したか
}

// revive:enable:exported
