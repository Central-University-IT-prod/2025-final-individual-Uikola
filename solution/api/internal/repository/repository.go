package repository

import (
	"api/internal/entity"
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// ClientRepository defines the interface for database operations related to clients.
type ClientRepository interface {
	// CreateMany inserts multiple client records into the database.
	// If a client with the same ID already exists, it should update the existing record.
	CreateMany(ctx context.Context, clients []entity.Client) ([]entity.Client, error)

	// GetByID retrieves a client from the database by its client_id.
	GetByID(ctx context.Context, clientID string) (entity.Client, error)
}

// AdvertiserRepository defines the interface for database operations related to advertisers.
type AdvertiserRepository interface {
	// CreateMany inserts multiple advertiser records into the database.
	// If an advertiser with the same primary key (AdvertiserID) already exists, it updates the existing record.
	CreateMany(ctx context.Context, advertisers []entity.Advertiser) ([]entity.Advertiser, error)

	// GetByID retrieves an advertiser by its AdvertiserID from the database.
	GetByID(ctx context.Context, advertiserID string) (entity.Advertiser, error)
}

// MLScoreRepository defines the interface for database operations related to ml-scores.
type MLScoreRepository interface {
	// Create inserts ml-score record into the database.
	// If a ml-score with the same ClientID and same AdvertiserID already exists, it updates the existing record.
	Create(ctx context.Context, mlScore entity.MLScore) error

	// Get retrieves a ml-score by its ClientID and AdvertiserID from the database.
	Get(ctx context.Context, clientID, advertiserID string) (entity.MLScore, error)

	AvgSTDDEVWithTargeting(ctx context.Context, clientID string, advertisersIDs []string) (float64, float64, error)
}

// CampaignRepository defines the interface for database operations related to campaigns.
type CampaignRepository interface {
	// Create inserts campaign record into the database.
	// If a campaign with the same ClientID and same AdvertiserID already exists, it updates the existing record.
	Create(ctx context.Context, campaign entity.Campaign) (entity.Campaign, error)

	// Get retrieves a campaign by its AdvertiserID and CampaignID from the database.
	Get(ctx context.Context, advertiserID, campaignID string) (entity.Campaign, error)

	// GetByID retrieves a campaign by its CampaignID from the database.
	GetByID(ctx context.Context, campaignID string) (entity.Campaign, error)

	// ListWithPagination retrieves a list of campaigns by theirs AdvertiserID from database with pagination.
	ListWithPagination(ctx context.Context, advertiserID string, limit, offset int) ([]entity.Campaign, error)

	// Count is a function that gets the count of advertiser's campaigns from the database.
	Count(ctx context.Context, advertiserID string) (int64, error)

	// ListForModerationWithPagination retrieves a list of campaigns for moderation from database with pagination.
	ListForModerationWithPagination(ctx context.Context, limit, offset int) ([]entity.Campaign, error)

	// CountForModeration is a function that gets the count of campaigns for moderation from the database.
	CountForModeration(ctx context.Context) (int64, error)

	// ListWithTargeting retrieves a list of active campaigns that match the given targeting criteria (gender, age, location).
	// It filters campaigns based on the targeting settings and ensures they are currently active.
	ListWithTargeting(ctx context.Context, client entity.Client, currentDay int) ([]entity.Campaign, error)

	GetAvgSTDDEVMLScorePerCampaignWithTargeting(ctx context.Context) (map[string]float64, map[string]float64, error)

	AvgSTDDEVCPICPC(ctx context.Context, client entity.Client, currentDay int) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error)

	// Update updates an existing campaign along with its associated targeting data within a transaction.
	Update(ctx context.Context, campaign entity.Campaign) error

	// Delete removes a campaign from the database by its AdvertiserID and CampaignID.
	Delete(ctx context.Context, advertiserID, campaignID string) error
}

type AdRepository interface {
	// TryIncrementImpressions attempts to increment the impressions count for a campaign.
	// Returns true if incremented successfully (within the limit), otherwise false
	TryIncrementImpressions(ctx context.Context, campaignID string, limit int) (bool, error)

	// IncrementClick increments the clicks count for a given campaign.
	IncrementClick(ctx context.Context, campaignID string) error

	// GetImpressionsCount retrieves the impressions count for a campaign.
	// Returns 0 if the key is not found.
	GetImpressionsCount(ctx context.Context, campaignID string) (int, error)

	// GetClicksCount retrieves the clicks count for a campaign.
	// Returns 0 if the key is not found.
	GetClicksCount(ctx context.Context, campaignID string) (int, error)

	// AddSeenCampaign records that a client has seen a specific campaign.
	AddSeenCampaign(ctx context.Context, clientID, campaignID string) error

	// HasSeenCampaign checks if a client has already seen a specific campaign.
	HasSeenCampaign(ctx context.Context, clientID, campaignID string) (bool, error)

	// AddClickedCampaign records that a client has clicked on a specific campaign.
	AddClickedCampaign(ctx context.Context, clientID, campaignID string) error

	// HasClickedCampaign checks if a client has already clicked on a specific campaign.
	HasClickedCampaign(ctx context.Context, clientID, campaignID string) (bool, error)
}

// TimeRepository defines the interface for database operations related to time.
type TimeRepository interface {
	// Get retrieves the stored day value from Redis.
	Get(ctx context.Context) (int, error)

	// Set stores a day value in Redis with an expiration time.
	Set(ctx context.Context, day int, expiration time.Duration)
}

// ImpressionRepository defines the interface for database operations related to impression.
type ImpressionRepository interface {
	// Create inserts impression record into the database.
	Create(ctx context.Context, impression entity.Impression) (entity.Impression, error)

	Get(ctx context.Context, advertiserID, campaignID, clientID string) (entity.Impression, error)
}

// ClickRepository defines the interface for database operations related to click.
type ClickRepository interface {
	// Create inserts click record into the database.
	Create(ctx context.Context, click entity.Click) (entity.Click, error)

	// Get retrieves a click by its AdvertiserID, CampaignID and ClientID from the database.
	Get(ctx context.Context, advertiserID, campaignID, clientID string) (entity.Click, error)
}

// StatisticRepository defines the interface for database operations related to statistic.
type StatisticRepository interface {
	// Delete removes records of impressions and clicks associated with the given advertiser and campaign.
	// The operation is performed within a single database transaction to ensure consistency.
	Delete(ctx context.Context, advertiserID, campaignID string) error

	// GetByCampaignID retrieves the aggregated statistics for a specific campaign by its ID.
	GetByCampaignID(ctx context.Context, campaignID string) (entity.Statistic, error)

	// GetByAdvertiserID retrieves aggregated statistics for a specific advertiser.
	GetByAdvertiserID(ctx context.Context, advertiserID string) (entity.Statistic, error)

	// GetDailyByCampaignID retrieves daily statistics for a specific campaign.
	GetDailyByCampaignID(ctx context.Context, campaignID string) ([]entity.Statistic, error)

	// GetDailyByAdvertiserID retrieves daily statistics for a specific advertiser.
	GetDailyByAdvertiserID(ctx context.Context, advertiserID string) ([]entity.Statistic, error)
}
