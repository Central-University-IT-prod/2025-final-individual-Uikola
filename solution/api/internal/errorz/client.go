package errorz

import "errors"

var (
	ErrInvalidClientID = errors.New("invalid client id")

	ErrClientNotFound = errors.New("client not found")
)
