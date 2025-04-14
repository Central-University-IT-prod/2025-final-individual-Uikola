package clickhouse

import (
	"api/internal/entity"
	"context"

	"gorm.io/gorm"
)

// StatisticRepository handles database operations for the Statistic entity.
type StatisticRepository struct {
	db *gorm.DB
}

// NewStatisticRepository creates a new instance of StatisticRepository.
func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{
		db: db,
	}
}

// Delete removes records of impressions and clicks associated with the given advertiser and campaign.
// The operation is performed within a single database transaction to ensure consistency.
func (r *StatisticRepository) Delete(ctx context.Context, advertiserID, campaignID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("advertiser_id = ? AND campaign_id = ?", advertiserID, campaignID).Delete(&entity.Impression{}).Error; err != nil {
			return err
		}
		if err := tx.Where("advertiser_id = ? AND campaign_id = ?", advertiserID, campaignID).Delete(&entity.Click{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByCampaignID retrieves the aggregated statistics for a specific campaign by its ID.
func (r *StatisticRepository) GetByCampaignID(ctx context.Context, campaignID string) (entity.Statistic, error) {
	var statistic entity.Statistic
	err := r.db.WithContext(ctx).
		Table("impressions i").
		Select(`
			COUNT(i.client_id) AS impressions_count,
			COUNT(CASE WHEN c.client_id != '00000000-0000-0000-0000-000000000000' THEN c.client_id END) AS clicks_count,
			SUM(i.cost_per_impression) AS spent_impressions,
			SUM(c.cost_per_click) AS spent_clicks
		`).
		Joins("LEFT JOIN clicks c ON i.advertiser_id = c.advertiser_id AND i.campaign_id = c.campaign_id AND i.client_id = c.client_id").
		Where("i.campaign_id = ?", campaignID).
		Scan(&statistic).Error

	return statistic, err
}

// GetByAdvertiserID retrieves aggregated statistics for a specific advertiser.
func (r *StatisticRepository) GetByAdvertiserID(ctx context.Context, advertiserID string) (entity.Statistic, error) {
	var statistic entity.Statistic
	err := r.db.WithContext(ctx).
		Table("impressions i").
		Select(`
			COUNT(i.client_id) AS impressions_count,
			COUNT(CASE WHEN c.client_id != '00000000-0000-0000-0000-000000000000' THEN c.client_id END) AS clicks_count,
			SUM(i.cost_per_impression) AS spent_impressions,
			SUM(c.cost_per_click) AS spent_clicks
		`).
		Joins("LEFT JOIN clicks c ON i.advertiser_id = c.advertiser_id AND i.campaign_id = c.campaign_id AND i.client_id = c.client_id").
		Where("i.advertiser_id = ?", advertiserID).
		Scan(&statistic).Error

	return statistic, err
}

// GetDailyByCampaignID retrieves daily statistics for a specific campaign.
func (r *StatisticRepository) GetDailyByCampaignID(ctx context.Context, campaignID string) ([]entity.Statistic, error) {
	var statistics []entity.Statistic
	err := r.db.WithContext(ctx).Raw(
		`SELECT date,
       				SUM(impressions_count) AS impressions_count,
       				SUM(clicks_count) AS clicks_count,
       				SUM(spent_impressions) AS spent_impressions,
       				SUM(spent_clicks) AS spent_clicks

	   		 FROM (
         			SELECT created_at AS date, COUNT(*) AS impressions_count, 0 AS clicks_count, SUM(cost_per_impression) AS spent_impressions, 0 AS spent_clicks
         			FROM impressions
         			WHERE campaign_id = ?
         			GROUP BY date

         			UNION ALL

         			SELECT created_at AS date, 0 AS impressions_count, COUNT(*) AS clicks_count, 0 AS spent_impressions, SUM(cost_per_click) AS spent_clicks
         			FROM clicks
         			WHERE campaign_id = ?
         			GROUP BY date
         ) AS combined
			 GROUP BY date
			 ORDER BY date`, campaignID, campaignID).Scan(&statistics).Error

	return statistics, err
}

// GetDailyByAdvertiserID retrieves daily statistics for a specific advertiser.
func (r *StatisticRepository) GetDailyByAdvertiserID(ctx context.Context, advertiserID string) ([]entity.Statistic, error) {
	var statistics []entity.Statistic
	err := r.db.WithContext(ctx).Raw(
		`SELECT date,
       				SUM(impressions_count) AS impressions_count,
       				SUM(clicks_count) AS clicks_count,
       				SUM(spent_impressions) AS spent_impressions,
       				SUM(spent_clicks) AS spent_clicks

	   		 FROM (
         			SELECT created_at AS date, COUNT(*) AS impressions_count, 0 AS clicks_count, SUM(cost_per_impression) AS spent_impressions, 0 AS spent_clicks
         			FROM impressions
         			WHERE advertiser_id = ?
         			GROUP BY date

         			UNION ALL

         			SELECT created_at AS date, 0 AS impressions_count, COUNT(*) AS clicks_count, 0 AS spent_impressions, SUM(cost_per_click) AS spent_clicks
         			FROM clicks
         			WHERE advertiser_id = ?
         			GROUP BY date
         ) AS combined
			 GROUP BY date
			 ORDER BY date`, advertiserID, advertiserID).Scan(&statistics).Error

	return statistics, err
}
