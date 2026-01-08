// Package config provides application configuration settings.
package config

import (
	"os"
)

// MySQLSettings holds the configuration settings for connecting to a MySQL database.
type MySQLSettings struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

// FirebaseSettings は Firebase 連携に必要な設定値です。
type FirebaseSettings struct {
	ProjectID        string
	CredsFile        string
	UseAuthEmulator  bool
	AuthEmulatorHost string
}

// Config はアプリケーションの設定を保持する構造体です。
type Config struct {
	MySQL      MySQLSettings
	AppPort    string
	AppIP      string
	AppLogPath string
	SQLLogPath string
	Firebase   FirebaseSettings
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Load は環境変数からアプリ設定を構築して返します。
func Load() Config {
	return Config{
		MySQL: MySQLSettings{
			Host:   getEnv("DB_HOST", "127.0.0.1"),
			Port:   getEnv("DB_PORT", "3306"),
			User:   getEnv("DB_USER", "root"),
			Pass:   getEnv("DB_PASS", "root"),
			DBName: getEnv("DB_NAME", "mydb"),
		},
		AppPort:    getEnv("APP_PORT", "8080"),
		AppIP:      getEnv("APP_IP", "127.0.0.1"),
		AppLogPath: getEnv("APP_LOG_PATH", "app.log"),
		SQLLogPath: getEnv("SQL_LOG_PATH", "app.log"),
		Firebase: FirebaseSettings{
			ProjectID:        getEnv("FIREBASE_PROJECT_ID", ""),
			CredsFile:        getEnv("FIREBASE_CREDENTIALS_FILE", ""),
			UseAuthEmulator:  getEnv("FIREBASE_AUTH_EMULATOR", "false") == "true",
			AuthEmulatorHost: getEnv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099"),
		},
	}
}
