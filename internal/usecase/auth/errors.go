// Package auth contains auth-usecase-specific thin adapters around shared errors.
// Keep this file only if you want to minimize diffs in existing call sites.
package auth

import (
	ucshared "resume/internal/usecase/shared"
)

// ErrInvalidInput legacy alias preserved for compatibility.
// Prefer using ucshared.ErrUnprocessable (422) or ucshared.ErrBadRequest (400) explicitly.
var ErrInvalidInput = ucshared.ErrUnprocessable

// mapInfraToUCError delegates to the shared mapper.
// Existing call sites in auth usecase can remain unchanged.
func mapInfraToUCError(err error) error {
	if err == nil {
		return nil
	}
	return ucshared.MapInfraError(err)
}
