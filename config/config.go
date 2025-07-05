// Package config provides application configuration settings.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// MySQLSettings holds the configuration settings for connecting to a MySQL database.
type MySQLSettings struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

// DSN はMySQL接続用のDSN文字列を返します。
func (c *MySQLSettings) DSN() string {
	return c.User + ":" + c.Pass + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.DBName + "?parseTime=True"
}

// Config はアプリケーションの設定を保持する構造体です。
type Config struct {
	MySQL      MySQLSettings
	AppPort    string
	AppIP      string
	AppLogPath string
	SQLLogPath string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}
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
	}
}
