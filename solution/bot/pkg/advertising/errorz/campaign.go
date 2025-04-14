package errorz

import "errors"

var (
	ErrCampaignNotFound = errors.New("campaign not found")
)
