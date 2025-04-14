package errorz

import "errors"

var (
	ErrInvalidData = errors.New("invalid data")
	ErrUnexpected  = errors.New("unexpected error")
)
