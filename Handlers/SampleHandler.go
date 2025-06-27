package Handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPing(c *gin.Context) {
	fmt.Println("test")
	//logger := slog.FromContext(c.Request.Context())

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
