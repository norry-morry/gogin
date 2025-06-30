package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTest(c *gin.Context) {
	fmt.Println("GetTest")
	c.JSON(http.StatusOK, gin.H{
		"message": "GetTest",
	})
}

func GetTest2(c *gin.Context) {
	fmt.Println("GetTest2")
	c.JSON(http.StatusOK, gin.H{
		"message": "GetTest2",
	})
}
