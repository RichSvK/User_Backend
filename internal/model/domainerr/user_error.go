package domainerr

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("invalid user credentials")
	ErrNotVerified   = errors.New("email not verified")
	ErrEmailExists   = errors.New("email already registered")
	ErrInternal      = errors.New("internal server error")
	ErrVerified      = errors.New("user is already verified")
	ErrUserLoggedIn  = errors.New("you are already logged in")

	// JWT related errors
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrMissingSubject     = errors.New("missing subject claim")
	ErrEmptyToken         = errors.New("token is required")
)
