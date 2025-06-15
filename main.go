package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"resume/Routes"
	"resume/library/database"
)

func main() {
	getEnv := godotenv.Load()
	if getEnv != nil {
		fmt.Println("Error loading .env file")
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	//r := gin.Default()
	r := Routes.SetupRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello air! pong2",
		})
	})
	err := r.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
