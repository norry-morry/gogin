// Package firebase provides factories and helpers to initialize
// the Firebase Admin SDK (Auth) for the application.
package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Credentials Config から必要な情報だけ受け取る形にしておくと依存が軽い
type Credentials struct {
	// JSON資格情報のファイルパス or そのままJSON文字列
	CredsFile string // ex) "/work/keys/service-account.json"
	CredsJSON []byte // ex) 埋め込みやSecretから渡す場合
	ProjectID string // GCP プロジェクトID（無くても動くが明示推奨）
}

// Options は Firebase クライアントの動作設定を保持する構造体です。
// 実環境またはエミュレータ環境のいずれを使用するかを指定します。
type Options struct {
	// UseAuthEmulator は Firebase Authentication エミュレータを利用する場合に true に設定します。
	// false の場合は実際の Firebase プロジェクトに接続します。
	UseAuthEmulator bool

	// AuthEmulatorHost は Authentication エミュレータのホスト名またはアドレスを指定します。
	// 例: "localhost:9099"
	AuthEmulatorHost string
}

// NewAuthClient は Firebase Admin の *auth.Client を生成して返す
func NewAuthClient(ctx context.Context, creds Credentials, opt Options) (*auth.Client, error) {
	var appOpts []option.ClientOption

	// 認証方法: ファイル or JSON（どちらか一方でOK）
	switch {
	case len(creds.CredsJSON) > 0:
		appOpts = append(appOpts, option.WithCredentialsJSON(creds.CredsJSON))
	case creds.CredsFile != "":
		appOpts = append(appOpts, option.WithCredentialsFile(creds.CredsFile))
	default:
		// GCE/GKE 等なら ADC (Application Default Credentials) にフォールバック
	}

	cfg := &firebase.Config{}
	if creds.ProjectID != "" {
		cfg.ProjectID = creds.ProjectID
	}

	// Emulator 対応（Admin SDKは環境変数を見る）
	if opt.UseAuthEmulator {
		host := opt.AuthEmulatorHost
		if host == "" {
			host = "localhost:9099"
		}
		// Admin SDK は FIREBASE_AUTH_EMULATOR_HOST を参照する。
		// セット失敗は致命ではない（後段の初期化で検出できる）ため握りつぶす。
		//_ = os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", host)
		if err := os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", host); err != nil {
			// ここで logger が使えない前提なので stderr に出す
			if _, e := fmt.Fprintf(os.Stderr,
				"warn: failed to set FIREBASE_AUTH_EMULATOR_HOST=%s: %v\n",
				host, err,
			); e != nil {
				// stderr への出力すら失敗。ここでは致命ではないので変数を参照して「扱った」ことにする。
				// （staticcheck SA9003 対策）
				_ = e
			}
		}
	} else {
		// Emulator を使わないときは強制的に無効化
		if err := os.Unsetenv("FIREBASE_AUTH_EMULATOR_HOST"); err != nil {
			if _, e := fmt.Fprintf(os.Stderr,
				"warn: failed to unset FIREBASE_AUTH_EMULATOR_HOST: %v\n", err,
			); e != nil {
				_ = e // staticcheck対策
			}
		}
	}

	app, err := firebase.NewApp(ctx, cfg, appOpts...)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %w", err)
	}

	return app.Auth(ctx)
}
