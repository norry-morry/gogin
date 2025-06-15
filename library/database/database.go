package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// DB データベース接続用のグローバル変数
var DB *gorm.DB

// DBConfig データベース接続情報の構造体
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// 環境変数からDB設定を取得
func getDBConfig() *DBConfig {
	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// DSN文字列を生成
func (c *DBConfig) buildDSN() string {
	//dsn := "root@tcp(mysql:3306)/local_db?charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}

// InitDB データベース初期化
func InitDB() {
	config := getDBConfig()
	dsn := config.buildDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true, // エラー翻訳機能を有効化
	})
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}

	DB = db
	log.Println("データベース接続に成功しました")
}
