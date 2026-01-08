// Package handlerutil は、Gin ベースの HTTP ハンドラで共通的に利用される
// エラーレスポンス整形や Request ID 解決などのユーティリティを提供します。
//
// 主な機能:
//   - WriteError: アプリケーション層のエラー（apperr.Error など）を
//     統一された JSON レスポンス形式に変換して返します。
//
// これにより、ハンドラごとの重複実装を避けつつ、
// 一貫したエラーレスポンス構造とトレーサビリティを確保できます。
package handlerutil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v10 "github.com/go-playground/validator/v10"

	verrs "resume/internal/adapter/validation/errors" // ToDetails(err) を置いた所
	"resume/internal/shared/apperr"
	"resume/internal/shared/requestid"
)

// レスポンス形（error をトップにネストしている前提）
type errorBody struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"requestId,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}
type errorEnvelope struct {
	Error errorBody `json:"error"`
}

// apperr.Code → HTTP ステータス（あなたの環境に HTTPStatusOf があるならそちらを使ってOK）
func httpStatusOf(code apperr.Code) int {
	return apperr.HTTPStatusOf(code)
}

// apperr.Code を安全に文字列化（String 実装があればそれが使われる）
func codeToString(code apperr.Code) string { return fmt.Sprint(code) }

// WriteError ハンドラ側は必ず `return` で呼び出しを終了してください（二重書き込み防止）
func WriteError(c *gin.Context, err error) {
	// デフォルト（予期せぬエラー）
	status := http.StatusInternalServerError
	body := errorEnvelope{
		Error: errorBody{
			Code:    "internal_error",
			Message: "internal server error",
			Details: nil,
		},
	}

	// エラーの種類を判定して Details などだけ差し替える
	var ae *apperr.Error
	var verr v10.ValidationErrors

	switch {
	case errors.As(err, &ae):
		status = httpStatusOf(ae.Code)
		body.Error.Code = codeToString(ae.Code)
		body.Error.Message = ae.Message
		body.Error.Details = ae.Details

	case errors.As(err, &verr):
		status = http.StatusUnprocessableEntity
		body.Error.Code = "UNPROCESSABLE"
		body.Error.Message = "The request contains semantically invalid data."
		body.Error.Details = verrs.ToDetails(err)

	default:
		// 必要なら err.Error() を details に入れてもOK（内部情報を出したくないなら nil のまま）
		// body.Error.Details = map[string]any{"error": err.Error()}
	}

	// requestId を最終段で埋める（あなたの既存コードに合わせて body.Error.RequestID にセット）
	rid := requestid.Get(c)
	body.Error.RequestID = rid
	// 返却ヘッダにも反映したい場合（任意）
	if rid != "" && c.Writer.Header().Get("X-Request-Id") == "" {
		c.Writer.Header().Set("X-Request-Id", rid)
	}

	// ★ 最後は必ず c.JSON で統一（Abort は使わない）
	c.JSON(status, body)
}
