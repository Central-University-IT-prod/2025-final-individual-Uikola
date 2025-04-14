package postgres

import (
	"api/internal/entity"
	"api/internal/errorz"
	"context"
	"errors"

	"gorm.io/gorm"
)

// AdvertiserRepository handles database operations for the Advertiser entity.
type AdvertiserRepository struct {
	db *gorm.DB
}

// NewAdvertiserRepository creates a new instance of AdvertiserRepository.
func NewAdvertiserRepository(db *gorm.DB) *AdvertiserRepository {
	return &AdvertiserRepository{
		db: db,
	}
}

// CreateMany inserts multiple advertiser records into the database.
// If an advertiser with the same primary key (AdvertiserID) already exists, it updates the existing record.
func (r *AdvertiserRepository) CreateMany(ctx context.Context, advertisers []entity.Advertiser) ([]entity.Advertiser, error) {
	err := r.db.WithContext(ctx).Save(&advertisers).Error
	return advertisers, err
}

// GetByID retrieves an advertiser by its AdvertiserID from the database.
func (r *AdvertiserRepository) GetByID(ctx context.Context, advertiserID string) (entity.Advertiser, error) {
	var advertiser entity.Advertiser
	err := r.db.WithContext(ctx).Where("advertiser_id = ?", advertiserID).First(&advertiser).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Advertiser{}, errorz.ErrAdvertiserNotFound
	case err != nil:
		return entity.Advertiser{}, err
	}

	return advertiser, nil
}
