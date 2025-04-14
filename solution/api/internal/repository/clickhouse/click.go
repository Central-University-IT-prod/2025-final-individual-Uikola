package clickhouse

import (
	"api/internal/entity"
	"api/internal/errorz"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ClickRepository handles database operations for the Click entity.
type ClickRepository struct {
	db *gorm.DB
}

// NewClickRepository creates a new instance of ClickRepository.
func NewClickRepository(db *gorm.DB) *ClickRepository {
	return &ClickRepository{
		db: db,
	}
}

// Create inserts click record into the database.
func (r *ClickRepository) Create(ctx context.Context, click entity.Click) (entity.Click, error) {
	err := r.db.WithContext(ctx).Create(&click).Error
	return click, err
}

// Get retrieves a click by its AdvertiserID, CampaignID and ClientID from the database.
func (r *ClickRepository) Get(ctx context.Context, advertiserID, campaignID, clientID string) (entity.Click, error) {
	var click entity.Click
	err := r.db.WithContext(ctx).Where("advertiser_id = ? AND campaign_id = ? AND client_id = ?", advertiserID, campaignID, clientID).First(&click).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Click{}, errorz.ErrClickNotFound
	case err != nil:
		return entity.Click{}, err
	}
	return click, err
}
