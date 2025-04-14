package clickhouse

import (
	"api/internal/entity"
	"api/internal/errorz"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ImpressionRepository handles database operations for the Impression entity.
type ImpressionRepository struct {
	db *gorm.DB
}

// NewImpressionRepository creates a new instance of ImpressionRepository.
func NewImpressionRepository(db *gorm.DB) *ImpressionRepository {
	return &ImpressionRepository{
		db: db,
	}
}

// Create inserts impression record into the database.
func (r *ImpressionRepository) Create(ctx context.Context, impression entity.Impression) (entity.Impression, error) {
	err := r.db.WithContext(ctx).Create(&impression).Error
	return impression, err
}

// Get retrieves an impression by its AdvertiserID, CampaignID and ClientID from the database.
func (r *ImpressionRepository) Get(ctx context.Context, advertiserID, campaignID, clientID string) (entity.Impression, error) {
	var impression entity.Impression
	err := r.db.WithContext(ctx).Where("advertiser_id = ? AND campaign_id = ? AND client_id = ?", advertiserID, campaignID, clientID).First(&impression).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Impression{}, errorz.ErrImpressionNotFound
	case err != nil:
		return entity.Impression{}, err
	}

	return impression, err

}
