// Package db Package config provides database connection configuration settings.
package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"resume/internal/config"
)

// NewGorm は GORM の db 接続を初期化して返します。
func NewGorm(cfg config.Config, gl logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.MySQL.User, cfg.MySQL.Pass, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.DBName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gl,
	})
}
