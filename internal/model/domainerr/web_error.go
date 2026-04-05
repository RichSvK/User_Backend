package domainerr

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrTooManyRequest     = errors.New("too many requests")

	ErrAuthorizationHeaderRequired = errors.New("authorization header is required")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrTokenExpired                = errors.New("token has expired")
	ErrUnauthorizedAccess          = errors.New("unauthorized access")
	ErrServiceTimeout              = errors.New("request timeout")
)
