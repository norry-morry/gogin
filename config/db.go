// Package config provides database connection configuration settings.
package config

import (
	"fmt"
	"log/slog"
	"resume/infrastructure/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewDB は指定されたMySQL設定とロガーを使用して新しいGormのDBインスタンスを返します。
func NewDB(cfg MySQLSettings, sqlLogger *slog.Logger) *gorm.DB {
	dsn := cfg.DSN()

	slogLogger := logger.NewSlogGormLogger(
		sqlLogger,
		gormlogger.Info,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: slogLogger,
	})
	if err != nil {
		slog.Default().Error("DB connection failed", slog.Any("error", err))
		panic(fmt.Sprintf("failed to connect to DB: %v", err))
		//log.Fatalf("failed to connect database: %v", err)
	}

	// DB接続情報(ping)
	sqlDB, err := db.DB()
	if err != nil {
		slog.Default().Error("failed to get generic DB", slog.Any("error", err))
		panic(err)
	}
	if err := sqlDB.Ping(); err != nil {
		slog.Default().Error("failed to ping DB", slog.Any("error", err))
		panic(err)
	}

	// 接続プール設定（任意）
	//sqlDB.SetMaxOpenConns(25)
	//sqlDB.SetMaxIdleConns(25)
	//sqlDB.SetConnMaxLifetime(5 * time.Minute)

	slog.Default().Info("successfully connected to DB")
	return db
}
