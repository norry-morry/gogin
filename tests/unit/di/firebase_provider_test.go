// Package di_test は internal/di パッケージで定義された DI プロバイダ群に対する
// ユニットテストを提供します。このファイルでは Firebase 認証関連の
// プロバイダのテストを行います。
package di_test

import (
	"testing"

	"resume/internal/config"
	"resume/internal/di"
)

// TestProvideFBCreds_MapsFromConfig は ExportProvideFBCreds が
// Firebase 設定を Credentials にマッピングすることを確認します。
func TestProvideFBCreds_MapsFromConfig(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Firebase: config.FirebaseSettings{
			ProjectID: "proj-123",
			CredsFile: "/secrets/firebase.json",
		},
	}

	creds := di.ExportProvideFBCreds(cfg)
	if creds.ProjectID != "proj-123" {
		t.Errorf("ProjectID = %s, want proj-123", creds.ProjectID)
	}
	if creds.CredsFile != "/secrets/firebase.json" {
		t.Errorf("CredsFile = %s, want /secrets/firebase.json", creds.CredsFile)
	}
}

// TestProvideFBOptions_MapsFromConfig は ExportProvideFBOptions が
// Firebase 設定を Options にマッピングすることを確認します。
func TestProvideFBOptions_MapsFromConfig(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Firebase: config.FirebaseSettings{
			UseAuthEmulator:  true,
			AuthEmulatorHost: "localhost:9099",
		},
	}

	opt := di.ExportProvideFBOptions(cfg)
	if !opt.UseAuthEmulator {
		t.Errorf("UseAuthEmulator = %v, want true", opt.UseAuthEmulator)
	}
	if opt.AuthEmulatorHost != "localhost:9099" {
		t.Errorf("AuthEmulatorHost = %s, want localhost:9099", opt.AuthEmulatorHost)
	}
}
