// Package handler はエンドポイント処理を担当します。
package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPing(c *gin.Context) {
	slog.Debug("GetPing called")
	//logger := slog.FromContext(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func GetSample(c *gin.Context) {
	slog.Debug("GetSample called")
	c.JSON(http.StatusOK, gin.H{
		"message": "sample",
	})
}
