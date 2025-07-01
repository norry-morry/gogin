package config

import (
	"github.com/joho/godotenv"
	"os"
)

type MySQLSettings struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

func (c *MySQLSettings) DSN() string {
	return c.User + ":" + c.Pass + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.DBName + "?parseTime=True"
}

type Config struct {
	MySQL      MySQLSettings
	AppPort    string
	AppIP      string
	AppLogPath string
	SqlLogPath string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	_ = godotenv.Load()
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
		SqlLogPath: getEnv("SQL_LOG_PATH", "app.log"),
	}
}
