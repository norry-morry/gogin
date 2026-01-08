// Package apperr は、アプリケーション全体で共通的に扱う
// エラーコードと HTTP ステータスマッピングを提供します。
package apperr

import "net/http"

// Code represents a stable, application-wide error code.
// It should be used for client-facing error classification and logging.
type Code string

const (
	// CodeBadRequest indicates malformed inputs: JSON decode errors,
	// missing required keys at transport-level, or type mismatches.
	CodeBadRequest Code = "BAD_REQUEST"

	// CodeUnprocessable indicates validation failures: structurally valid
	// inputs that fail business/field validation rules.
	CodeUnprocessable Code = "UNPROCESSABLE"

	// CodeNotFound indicates that a requested single resource does not exist.
	CodeNotFound Code = "NOT_FOUND"

	// CodeConflict indicates a write conflict such as unique constraint violations.
	CodeConflict Code = "CONFLICT"

	// CodeUnauthorized indicates authentication is required and has failed or not provided.
	CodeUnauthorized Code = "UNAUTHORIZED"

	// CodeForbidden indicates the authenticated principal lacks permission.
	CodeForbidden Code = "FORBIDDEN"

	// CodeInternal indicates an unexpected server-side error.
	CodeInternal Code = "INTERNAL"
)

// HTTPStatusOf returns a conventional HTTP status code for a given application Code.
func HTTPStatusOf(code Code) int {
	switch code {
	case CodeBadRequest:
		return http.StatusBadRequest // 400
	case CodeUnprocessable:
		return http.StatusUnprocessableEntity // 422
	case CodeUnauthorized:
		return http.StatusUnauthorized // 401
	case CodeForbidden:
		return http.StatusForbidden // 403
	case CodeNotFound:
		return http.StatusNotFound // 404
	case CodeConflict:
		return http.StatusConflict // 409
	default:
		return http.StatusInternalServerError // 500
	}
}
