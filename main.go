package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log/slog"
	"resume/Routes"
	"resume/logger"
)

func main() {
	getEnv := godotenv.Load()
	if getEnv != nil {
		fmt.Println("Error loading .env file")
	}

	// リクエストログを含む各種ロガーの設定
	log := logger.SetupLogger()
	slog.SetDefault(log)
	router := Routes.SetupRouter(log)

	// その他のミドルウェアを組込(使わないかもしれない？)
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		slog.Debug("PingHandler called") // default logger 使用
		c.JSON(200, gin.H{
			"message": "Hello air! pong2",
		})
	})
	fmt.Println("Hello golang from docker with air!")
	router.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
