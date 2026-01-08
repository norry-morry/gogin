// Package controller はエンドポイント処理を担当します。
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestHandler は疎通確認用のハンドラです。
type TestHandler struct{}

// NewTestHandler は TestHandler の新しいインスタンスを生成して返します。
// テストの HTTP リクエストを処理するためのハンドラを初期化します。
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// Ping は サンプルメッセージを返します。
func (h TestHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "test pong",
	})
}
