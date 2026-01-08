// Package controller はエンドポイント処理を担当します。
package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SampleHandler はサンプルのハンドラです
type SampleHandler struct{}

// NewSampleHandler は SampleHandler の新しいインスタンスを生成して返します
// サンプルのHTTPリクエストを処理するためのハンドラを初期化染ます
func NewSampleHandler() *SampleHandler { return &SampleHandler{} }

// Ping は サンプルメッセージを返します。
func (h SampleHandler) Ping(c *gin.Context) {
	slog.Debug("GetPing called")
	// logger := slog.FromContext(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"message": "sample pong",
	})
}
