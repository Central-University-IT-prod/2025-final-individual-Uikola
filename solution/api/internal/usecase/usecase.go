package usecase

import (
	"api/internal/entity/request"
	"api/internal/entity/response"
	"api/pkg/s3"
	"context"
)

// ClientUsecase defines the business logic layer for client operations.
type ClientUsecase interface {
	// CreateBulk processes and stores multiple client records.
	CreateBulk(ctx context.Context, req []request.Client) ([]response.Client, error)

	// GetByID retrieves a client by their client ID.
	GetByID(ctx context.Context, clientID string) (response.Client, error)
}

// AdvertiserUsecase defines the business logic layer for advertiser operations.
type AdvertiserUsecase interface {
	// CreateBulk processes and stores multiple advertisers records.
	CreateBulk(ctx context.Context, req []request.Advertiser) ([]response.Advertiser, error)

	// GetByID retrieves an advertiser by their advertiser ID.
	GetByID(ctx context.Context, advertiserID string) (response.Advertiser, error)
}

// MLScoreUsecase defines the business logic layer for ml-score operations.
type MLScoreUsecase interface {
	// Create processes a ml-score creation request.
	Create(ctx context.Context, req request.MLScore) error
}

// CampaignUsecase defines the business logic layer for campaign operations.
type CampaignUsecase interface {
	// Create processes a campaign creation request.
	Create(ctx context.Context, req request.CreateCampaign, advertiserID string) (response.Campaign, error)

	// Delete removes a campaign by its AdvertiserID and CampaignID.
	Delete(ctx context.Context, advertiserID, campaignID string) error

	// GetByID retrieves a campaign by its advertiser and campaign id.
	GetByID(ctx context.Context, advertiserID, campaignID string) (response.Campaign, error)

	// ListWithPagination retrieves a list of campaigns by its advertiser and campaign id with pagination.
	ListWithPagination(ctx context.Context, advertiserID string, size, page int) ([]response.Campaign, int64, error)

	// ListForModerationWithPagination retrieves a list of campaigns for moderation with pagination.
	ListForModerationWithPagination(ctx context.Context, size, page int) ([]response.Campaign, int64, error)

	GetAd(ctx context.Context, clientID string) (response.GetAd, error)

	// Update processes a campaign update request.
	Update(ctx context.Context, advertiserID, campaignID string, req request.UpdateCampaign) (response.Campaign, error)

	// ClickAd records a client's click on an advertisement.
	ClickAd(ctx context.Context, req request.ClickAd, adID string) error

	// UploadImage uploads an image to S3, updates the campaign with the image ID and returns the image URL.
	UploadImage(ctx context.Context, fileData s3.FileDataType, advertiserID, campaignID string) (string, error)

	// RemoveImage removes the image from the campaign by clearing the ImageID field.
	RemoveImage(ctx context.Context, advertiserID, campaignID string) error

	// GenerateText generates ad text based on the provided ad title and advertiser information.
	// It uses the AI client to create a compelling advertisement message.
	GenerateText(ctx context.Context, req request.GenerateAdText, advertiserID string) (response.GenerateAdText, error)

	// Moderate reviews the advertising campaign.
	Moderate(ctx context.Context, campaignID string, passedModeration bool) error
}

// StatisticUsecase defines the business logic layer for statistic operations.
type StatisticUsecase interface {
	// GetCampaign retrieves aggregated statistics for a specific campaign.
	GetCampaign(ctx context.Context, campaignID string) (response.GetStatistic, error)

	// GetAdvertiser retrieves the aggregated statistics for a specific advertiser.
	GetAdvertiser(ctx context.Context, advertiserID string) (response.GetStatistic, error)

	// GetCampaignDaily retrieves daily statistics for a specific campaign.
	GetCampaignDaily(ctx context.Context, campaignID string) ([]response.DailyStatistic, error)

	// GetAdvertiserDaily retrieves daily statistics for a specific advertiser.
	GetAdvertiserDaily(ctx context.Context, advertiserID string) ([]response.DailyStatistic, error)
}
