package errorz

import "errors"

var (
	ErrInvalidCampaignID = errors.New("invalid campaign id")
	ErrInvalidDate       = errors.New("invalid date")

	ErrCampaignNotFound  = errors.New("campaign not found")
	ErrCampaignIsActive  = errors.New("campaign is active")
	ErrCampaignIsOver    = errors.New("campaign is over")
	ErrNoCampaignsFound  = errors.New("no campaigns found")
	ErrCampaignIsNotSeen = errors.New("campaign is not seen")
)
