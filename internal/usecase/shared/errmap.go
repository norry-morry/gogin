// Package shared contains usecase-wide helpers for error normalization.
// It centralizes translation from infra-layer errors into application errors,
// and provides convenience helpers for common cases (e.g., not-found on nil).
package shared

import (
	"strings"

	"resume/internal/shared/apperr"
)

// ErrBadRequest should be used for transport-level issues detected in handlers
// (JSON decode failure, missing keys at payload level, type mismatches).
var ErrBadRequest = apperr.New(apperr.CodeBadRequest, "bad request", nil)

// ErrUnprocessable should be used for validation failures in handlers or usecases
// where the payload is syntactically valid but violates field/business rules.
var ErrUnprocessable = apperr.New(apperr.CodeUnprocessable, "validation failed", nil)

// MapInfraError normalizes infra-layer errors (db/driver/SDK) into application errors.
// NOTE: This version uses string heuristics as a stopgap.
// In production, prefer repository-defined sentinels and errors.Is checks.
func MapInfraError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "duplicate"), strings.Contains(msg, "unique"):
		return apperr.Wrap(apperr.CodeConflict, "resource already exists", err, nil)
	case strings.Contains(msg, "not found"), strings.Contains(msg, "record not found"):
		return apperr.Wrap(apperr.CodeNotFound, "resource not found", err, nil)
	default:
		return apperr.Wrap(apperr.CodeInternal, "internal error", err, nil)
	}
}

// NotFoundIfNil turns a nil pointer result into a NOT_FOUND application error.
// Useful for "findOne" style usecases to keep the calling code simple.
func NotFoundIfNil[T any](v *T, what string) (*T, error) {
	if v == nil {
		return nil, apperr.New(apperr.CodeNotFound, what+" not found", nil)
	}
	return v, nil
}
