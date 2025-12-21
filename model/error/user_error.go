package domain_error

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrNotVerified   = errors.New("email not verified")
	ErrEmailExists   = errors.New("email already registered")
	ErrInternal      = errors.New("internal server error")

	// JWT related errors
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrMissingSubject     = errors.New("missing subject claim")
)