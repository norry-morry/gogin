package main

import (
	"github.com/gin-gonic/gin"
	"resume/Routes"
)

func main() {
	//r := gin.Default()
	r := Routes.SetupRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello air! pong2",
		})
	})
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
