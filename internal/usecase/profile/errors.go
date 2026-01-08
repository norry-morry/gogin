package profile

import "resume/internal/shared/apperr"

func errUnauthorized() error {
	return apperr.New(apperr.CodeUnauthorized, "unauthorized", nil)
}

func errCheckAddressFailed(cause error) error {
	fields := map[string]any{}
	if cause != nil {
		fields["cause"] = cause.Error()
	}
	return apperr.New(apperr.CodeInternal, "failed to check address existence", fields)
}
