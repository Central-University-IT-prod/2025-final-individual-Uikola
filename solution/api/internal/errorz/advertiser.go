package errorz

import "errors"

var (
	ErrInvalidAdvertiserID = errors.New("invalid advertiser id")

	ErrAdvertiserNotFound = errors.New("advertiser not found")
)
