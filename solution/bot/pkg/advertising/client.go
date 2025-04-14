package advertising

import (
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
)

type Client interface {
	// GetClientByID retrieves a client by his ID.
	GetClientByID(clientID string) (response.Client, error)

	// GetAdvertiserByID retrieves an advertiser by his ID.
	GetAdvertiserByID(advertiserID string) (response.Advertiser, error)

	// GetCurrentDate retrieves a current api date.
	GetCurrentDate() (int, error)

	// CreateCampaign create campaign using provided request data.
	CreateCampaign(requestBody request.CreateCampaign, advertiserID string) (response.Campaign, error)

	// GetCampaign retrieves campaign with provided advertiser and campaign id.
	GetCampaign(advertiserID, campaignID string) (response.Campaign, error)

	// ListCampaigns retrieves a list of the advertiser campaigns.
	ListCampaigns(advertiserID string, size, page int) ([]response.Campaign, int, error)

	// UpdateCampaign update campaign with provided advertiser and campaign id.
	UpdateCampaign(requestBody request.UpdateCampaign, advertiserID, campaignID string) (response.Campaign, error)

	// DeleteCampaign delete campaign with provided advertiser and campaign id.
	DeleteCampaign(advertiserID, campaignID string) error

	// GetAdvertiserStats retrieves an advertiser stats by his ID.
	GetAdvertiserStats(advertiserID string) (response.Statistic, error)

	// GetCampaignStats retrieves a campaign stats by his ID.
	GetCampaignStats(campaignID string) (response.Statistic, error)

	GenerateAdText(requestBody request.GenerateAdText, advertiserID string) (response.GenerateAdText, error)

	GetAd(clientID string) (response.GetAd, error)

	ClickAd(requestBody request.ClickAd, campaignID string) error
}
