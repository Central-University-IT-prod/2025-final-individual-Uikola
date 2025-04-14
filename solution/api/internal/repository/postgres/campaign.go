package postgres

import (
	"api/internal/config"
	"api/internal/entity"
	"api/internal/errorz"
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// CampaignRepository handles database operations for the Campaign entity.
type CampaignRepository struct {
	db *gorm.DB

	moderationConfig config.ModerationConfig
}

// NewCampaignRepository creates a new instance of CampaignRepository.
func NewCampaignRepository(db *gorm.DB, moderationConfig config.ModerationConfig) *CampaignRepository {
	return &CampaignRepository{
		db: db,

		moderationConfig: moderationConfig,
	}
}

// Create inserts campaign record into the database.
// If a campaign with the same ClientID and same AdvertiserID already exists, it updates the existing record.
func (r *CampaignRepository) Create(ctx context.Context, campaign entity.Campaign) (entity.Campaign, error) {
	err := r.db.WithContext(ctx).Create(&campaign).Error
	return campaign, err
}

// Get retrieves a campaign by its AdvertiserID and CampaignID from the database.
func (r *CampaignRepository) Get(ctx context.Context, advertiserID, campaignID string) (entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.db.WithContext(ctx).Preload("Targeting").Where("advertiser_id = ? AND campaign_id = ?", advertiserID, campaignID).First(&campaign).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Campaign{}, errorz.ErrCampaignNotFound
	case err != nil:
		return entity.Campaign{}, err
	}

	return campaign, err
}

// GetByID retrieves a campaign by its CampaignID from the database.
func (r *CampaignRepository) GetByID(ctx context.Context, campaignID string) (entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.db.WithContext(ctx).Preload("Targeting").Where("campaign_id = ?", campaignID).First(&campaign).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Campaign{}, errorz.ErrCampaignNotFound
	case err != nil:
		return entity.Campaign{}, err
	}

	return campaign, err
}

// ListWithPagination retrieves a list of campaigns by theirs AdvertiserID from database with pagination.
func (r *CampaignRepository) ListWithPagination(ctx context.Context, advertiserID string, limit, offset int) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	err := r.db.WithContext(ctx).Preload("Targeting").Where("advertiser_id = ?", advertiserID).Limit(limit).Offset(offset).Find(&campaigns).Error
	return campaigns, err
}

// Count is a function that gets the count of advertiser's campaigns from the database.
func (r *CampaignRepository) Count(ctx context.Context, advertiserID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Campaign{}).Where("advertiser_id = ?", advertiserID).Count(&count).Error
	return count, err
}

// ListForModerationWithPagination retrieves a list of campaigns for moderation from database with pagination.
func (r *CampaignRepository) ListForModerationWithPagination(ctx context.Context, limit, offset int) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	err := r.db.WithContext(ctx).Preload("Targeting").Where("NOT passed_moderation").Limit(limit).Offset(offset).Find(&campaigns).Error
	return campaigns, err
}

// CountForModeration is a function that gets the count of campaigns for moderation from the database.
func (r *CampaignRepository) CountForModeration(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Campaign{}).Where("NOT passed_moderation").Count(&count).Error
	return count, err
}

// ListWithTargeting retrieves a list of active campaigns that match the given targeting criteria (gender, age, location).
// It filters campaigns based on the targeting settings, ensures they are currently active and have not been viewed by the client.
func (r *CampaignRepository) ListWithTargeting(ctx context.Context, client entity.Client, currentDay int) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	query := r.db.WithContext(ctx).
		Preload("Targeting").
		Joins("INNER JOIN targetings ON targetings.campaign_id = campaigns.campaign_id").
		Where("campaigns.start_date <= ? AND campaigns.end_date >= ?", currentDay, currentDay).
		Where("(targetings.gender IS NULL OR targetings.gender = ? OR targetings.gender = ?)", entity.TargetingGenderAll, client.Gender).
		Where("(targetings.age_from IS NULL OR targetings.age_from <= ?)", client.Age).
		Where("(targetings.age_to IS NULL OR targetings.age_to >= ?)", client.Age).
		Where("(targetings.location IS NULL OR targetings.location = ?)", client.Location)

	if r.moderationConfig.AdModerationEnabled() {
		query = query.Where("campaigns.passed_moderation")
	}

	err := query.Find(&campaigns).Error

	return campaigns, err
}

func (r *CampaignRepository) GetAvgSTDDEVMLScorePerCampaignWithTargeting(ctx context.Context) (map[string]float64, map[string]float64, error) {
	type campaignMLScore struct {
		CampaignID  string  `gorm:"column:campaign_id"`
		AvgScore    float64 `gorm:"column:avg_score"`
		STDDEVScore float64 `gorm:"column:stddev_score"`
	}
	var results []campaignMLScore

	err := r.db.WithContext(ctx).
		Table("campaigns").
		Select("campaigns.campaign_id, AVG(ml_scores.score) as avg_score, STDDEV(ml_scores.score) as stddev_score").
		Joins("INNER JOIN targetings ON campaigns.campaign_id = targetings.campaign_id").
		Joins("INNER JOIN ml_scores ON campaigns.advertiser_id = ml_scores.advertiser_id").
		Joins("INNER JOIN clients ON clients.client_id = ml_scores.client_id").
		Where("(targetings.gender IS NULL OR targetings.gender = 'ALL' OR targetings.gender = clients.gender)").
		Where("(targetings.age_from IS NULL OR targetings.age_from <= clients.age)").
		Where("(targetings.age_to IS NULL OR targetings.age_to >= clients.age)").
		Where("(targetings.location IS NULL OR targetings.location = clients.location)").
		Group("campaigns.campaign_id").
		Find(&results).Error

	if err != nil {
		return nil, nil, err
	}

	campaignScores := make(map[string]float64, len(results))
	stddevScores := make(map[string]float64, len(results))
	for _, result := range results {
		campaignScores[result.CampaignID] = result.AvgScore
		stddevScores[result.CampaignID] = result.STDDEVScore
	}

	return campaignScores, stddevScores, nil
}

func (r *CampaignRepository) AvgSTDDEVCPICPC(ctx context.Context, client entity.Client, currentDay int) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	var avgCPI, stddevCPI, avgCPC, stddevCPC decimal.Decimal
	query := r.db.WithContext(ctx).
		Preload("Targeting").
		Model(&entity.Campaign{}).
		Select("AVG(cost_per_impression), COALESCE(STDDEV(cost_per_impression), 0), AVG(cost_per_click), COALESCE(STDDEV(cost_per_click), 0)").
		Joins("INNER JOIN targetings ON targetings.campaign_id = campaigns.campaign_id").
		Where("campaigns.start_date <= ? AND campaigns.end_date >= ?", currentDay, currentDay).
		Where("(targetings.gender IS NULL OR targetings.gender = ? OR targetings.gender = ?)", entity.TargetingGenderAll, client.Gender).
		Where("(targetings.age_from IS NULL OR targetings.age_from <= ?)", client.Age).
		Where("(targetings.age_to IS NULL OR targetings.age_to >= ?)", client.Age).
		Where("(targetings.location IS NULL OR targetings.location = ?)", client.Location)

	if r.moderationConfig.AdModerationEnabled() {
		query = query.Where("campaigns.passed_moderation")
	}

	err := query.Row().Scan(&avgCPI, &stddevCPI, &avgCPC, &stddevCPC)

	return avgCPI, stddevCPI, avgCPC, stddevCPC, err
}

// Update updates an existing campaign along with its associated targeting data within a transaction.
func (r *CampaignRepository) Update(ctx context.Context, campaign entity.Campaign) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&campaign).Error; err != nil {
			return err
		}
		if err := tx.Save(&campaign.Targeting).Error; err != nil {
			return err
		}
		return nil
	})
}

// Delete removes a campaign from the database by its AdvertiserID and CampaignID.
func (r *CampaignRepository) Delete(ctx context.Context, advertiserID, campaignID string) error {
	err := r.db.WithContext(ctx).Where("advertiser_id = ? AND campaign_id = ?", advertiserID, campaignID).Delete(&entity.Campaign{}).Error
	return err
}
