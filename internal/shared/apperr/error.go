// Package apperr provides a small error type carrying an application Code,
// a client-safe message, optional details, and a wrapped cause.
// This lets usecases and handlers communicate stable error shapes to clients,
// while preserving the original error for logs and tracing via errors.Unwrap.
package apperr

import (
	"errors"
	"fmt"
)

// Error is the application error that flows through usecase→handler.
// - Code    : machine-friendly classification
// - Message : client-safe, human-readable text
// - Details : optional structured data (e.g., field errors, retry hints)
// - Cause   : wrapped underlying error for logging/diagnostics
type Error struct {
	Code    Code
	Message string
	Details map[string]any
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return string(e.Code)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
func (e *Error) Unwrap() error { return e.Cause }

// New constructs a new app error without a wrapped cause.
func New(code Code, msg string, details map[string]any) *Error {
	return &Error{Code: code, Message: msg, Details: details}
}

// Wrap constructs a new app error with a wrapped cause.
func Wrap(code Code, msg string, cause error, details map[string]any) *Error {
	return &Error{Code: code, Message: msg, Cause: cause, Details: details}
}

// Is reports whether err is an *apperr.Error with the given Code.
func Is(err error, code Code) bool {
	var ae *Error
	return errors.As(err, &ae) && ae.Code == code
}
