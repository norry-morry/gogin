package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPing(c *gin.Context) {
	fmt.Println("test")
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func GetSample(c *gin.Context) {
	fmt.Println("sample")
	c.JSON(http.StatusOK, gin.H{
		"message": "sample",
	})
}
