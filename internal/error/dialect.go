package errin

import "errors"

//goland:noinspection ALL
var (
	// ErrTokenRequired indicates no token given
	ErrTokenRequired  = errors.New("Token required")
	ErrMalformedToken = errors.New("Malformed token")
	ErrInvalidToken   = errors.New("Invalid token")

	// ErrInvalidID indicates invalid id
	ErrInvalidID = errors.New("invalid id format")

	// ErrInvalidStatus indicates invalid id
	ErrInvalidStatus = errors.New("invalid status")

	// ErrInvalidClaimContext indicates invalid action when claiming context
	ErrInvalidClaimContext = errors.New("unexpected error claim context")
)
