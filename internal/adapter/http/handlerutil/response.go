package handlerutil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	rules "resume/internal/adapter/validation/rules"
)

// ErrorItem は 1 件の検証エラー項目を表す型です。
// 検証タグ名・パラメータ・入力値などを保持します。
type ErrorItem = rules.ErrorItem

// ErrorEnvelope は 422 Unprocessable Entity 用のエラーレスポンス本体です。
// code / message / requestId と、フィールドごとの details を含みます。
type ErrorEnvelope struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"requestId,omitempty"`
	Details   map[string][]ErrorItem `json:"details,omitempty"`
	// 移行期だけ有効化したい場合はポインタでnil可
	LegacyDetails map[string]map[string]string `json:"legacyDetails,omitempty"`
}

// WriteValidationError は与えられた ErrorEnvelope を "error" キーでラップし、
// 422 Unprocessable Entity として JSON レスポンスを返します。
func WriteValidationError(c *gin.Context, env ErrorEnvelope) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"error": env})
}
