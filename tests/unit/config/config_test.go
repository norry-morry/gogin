// Package config_test は internal/config の Load() をテストします。
package config_test

import (
	"testing"

	"resume/internal/config"
)

//
// デフォルト値のテスト（環境変数が空の場合）
//

func TestLoad_DefaultsWhenEnvEmpty(t *testing.T) {
	// t.Setenv はテストごとに環境変数を上書き＆自動復元してくれる
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASS", "")
	t.Setenv("DB_NAME", "")

	t.Setenv("APP_PORT", "")
	t.Setenv("APP_IP", "")
	t.Setenv("APP_LOG_PATH", "")
	t.Setenv("SQL_LOG_PATH", "")

	t.Setenv("FIREBASE_PROJECT_ID", "")
	t.Setenv("FIREBASE_CREDENTIALS_FILE", "")
	t.Setenv("FIREBASE_AUTH_EMULATOR", "")
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "")

	got := config.Load()

	// MySQL
	if got.MySQL.Host != "127.0.0.1" {
		t.Errorf("MySQL.Host: expected 127.0.0.1, got %s", got.MySQL.Host)
	}
	if got.MySQL.Port != "3306" {
		t.Errorf("MySQL.Port: expected 3306, got %s", got.MySQL.Port)
	}
	if got.MySQL.User != "root" {
		t.Errorf("MySQL.User: expected root, got %s", got.MySQL.User)
	}
	if got.MySQL.Pass != "root" {
		t.Errorf("MySQL.Pass: expected root, got %s", got.MySQL.Pass)
	}
	if got.MySQL.DBName != "mydb" {
		t.Errorf("MySQL.DBName: expected mydb, got %s", got.MySQL.DBName)
	}

	// App
	if got.AppPort != "8080" {
		t.Errorf("AppPort: expected 8080, got %s", got.AppPort)
	}
	if got.AppIP != "127.0.0.1" {
		t.Errorf("AppIP: expected 127.0.0.1, got %s", got.AppIP)
	}
	if got.AppLogPath != "app.log" {
		t.Errorf("AppLogPath: expected app.log, got %s", got.AppLogPath)
	}
	if got.SQLLogPath != "app.log" {
		t.Errorf("SQLLogPath: expected app.log, got %s", got.SQLLogPath)
	}

	// Firebase
	if got.Firebase.ProjectID != "" {
		t.Errorf("Firebase.ProjectID: expected empty, got %s", got.Firebase.ProjectID)
	}
	if got.Firebase.CredsFile != "" {
		t.Errorf("Firebase.CredsFile: expected empty, got %s", got.Firebase.CredsFile)
	}
	if got.Firebase.UseAuthEmulator {
		t.Errorf("Firebase.UseAuthEmulator: expected false, got true")
	}
	if got.Firebase.AuthEmulatorHost != "localhost:9099" {
		t.Errorf("Firebase.AuthEmulatorHost: expected localhost:9099, got %s", got.Firebase.AuthEmulatorHost)
	}
}

//
// 環境変数で上書きできることのテスト
//

func TestLoad_EnvOverridesDefaults(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.local")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "appuser")
	t.Setenv("DB_PASS", "secret")
	t.Setenv("DB_NAME", "appdb")

	t.Setenv("APP_PORT", "9000")
	t.Setenv("APP_IP", "0.0.0.0")
	t.Setenv("APP_LOG_PATH", "/var/log/app.log")
	t.Setenv("SQL_LOG_PATH", "/var/log/sql.log")

	t.Setenv("FIREBASE_PROJECT_ID", "proj-123")
	t.Setenv("FIREBASE_CREDENTIALS_FILE", "/secrets/firebase.json")
	t.Setenv("FIREBASE_AUTH_EMULATOR", "true")
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9999")

	got := config.Load()

	// MySQL
	if got.MySQL.Host != "db.example.local" {
		t.Errorf("MySQL.Host: expected db.example.local, got %s", got.MySQL.Host)
	}
	if got.MySQL.Port != "3307" {
		t.Errorf("MySQL.Port: expected 3307, got %s", got.MySQL.Port)
	}
	if got.MySQL.User != "appuser" {
		t.Errorf("MySQL.User: expected appuser, got %s", got.MySQL.User)
	}
	if got.MySQL.Pass != "secret" {
		t.Errorf("MySQL.Pass: expected secret, got %s", got.MySQL.Pass)
	}
	if got.MySQL.DBName != "appdb" {
		t.Errorf("MySQL.DBName: expected appdb, got %s", got.MySQL.DBName)
	}

	// App
	if got.AppPort != "9000" {
		t.Errorf("AppPort: expected 9000, got %s", got.AppPort)
	}
	if got.AppIP != "0.0.0.0" {
		t.Errorf("AppIP: expected 0.0.0.0, got %s", got.AppIP)
	}
	if got.AppLogPath != "/var/log/app.log" {
		t.Errorf("AppLogPath: expected /var/log/app.log, got %s", got.AppLogPath)
	}
	if got.SQLLogPath != "/var/log/sql.log" {
		t.Errorf("SQLLogPath: expected /var/log/sql.log, got %s", got.SQLLogPath)
	}

	// Firebase
	if got.Firebase.ProjectID != "proj-123" {
		t.Errorf("Firebase.ProjectID: expected proj-123, got %s", got.Firebase.ProjectID)
	}
	if got.Firebase.CredsFile != "/secrets/firebase.json" {
		t.Errorf("Firebase.CredsFile: expected /secrets/firebase.json, got %s", got.Firebase.CredsFile)
	}
	if !got.Firebase.UseAuthEmulator {
		t.Errorf("Firebase.UseAuthEmulator: expected true, got false")
	}
	if got.Firebase.AuthEmulatorHost != "localhost:9999" {
		t.Errorf("Firebase.AuthEmulatorHost: expected localhost:9999, got %s", got.Firebase.AuthEmulatorHost)
	}
}
