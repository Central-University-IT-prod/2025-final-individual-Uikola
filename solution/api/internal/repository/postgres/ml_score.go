package postgres

import (
	"api/internal/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

// MLScoreRepository handles database operations for the MLScore entity.
type MLScoreRepository struct {
	db *gorm.DB
}

// NewMLScoreRepository creates a new instance of MLScoreRepository.
func NewMLScoreRepository(db *gorm.DB) *MLScoreRepository {
	return &MLScoreRepository{
		db: db,
	}
}

// Create inserts ml-score record into the database.
// If a ml-score with the same ClientID and same AdvertiserID already exists, it updates the existing record.
func (r *MLScoreRepository) Create(ctx context.Context, mlScore entity.MLScore) error {
	err := r.db.WithContext(ctx).Save(&mlScore).Error
	return err
}

// Get retrieves a ml-score by its ClientID and AdvertiserID from the database.
func (r *MLScoreRepository) Get(ctx context.Context, clientID, advertiserID string) (entity.MLScore, error) {
	var mlScore entity.MLScore
	err := r.db.WithContext(ctx).Where("client_id = ? AND advertiser_ID = ?", clientID, advertiserID).First(&mlScore).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.MLScore{
			ClientID:     clientID,
			AdvertiserID: advertiserID,
			Score:        0,
		}, nil
	case err != nil:
		return entity.MLScore{}, err
	}

	return mlScore, nil
}

func (r *MLScoreRepository) AvgSTDDEVWithTargeting(ctx context.Context, clientID string, advertisersIDs []string) (float64, float64, error) {
	var avgScore, stddev float64
	err := r.db.WithContext(ctx).
		Table("advertisers").
		Select("AVG(COALESCE(ml_scores.score, 0)), COALESCE(STDDEV(COALESCE(ml_scores.score, 0)), 0)").
		Joins("LEFT JOIN ml_scores ON advertisers.advertiser_id = ml_scores.advertiser_id AND ml_scores.client_id = ?", clientID).
		Where("advertisers.advertiser_id IN (?)", advertisersIDs).
		Row().Scan(&avgScore, &stddev)

	return avgScore, stddev, err
}

//func (r *MLScoreRepository) STDDEVWithTargeting(ctx context.Context, clientID string, advertisersIDs []string) (float64, error) {
//	var stddev float64
//	err := r.db.WithContext(ctx).
//		Table("advertisers").
//		Select("COALESCE(STDDEV(COALESCE(ml_scores.score, 0)), 0)").
//		Joins("LEFT JOIN ml_scores ON advertisers.advertiser_id = ml_scores.advertiser_id AND ml_scores.client_id = ?", clientID).
//		Where("advertisers.advertiser_id IN (?)", advertisersIDs).
//		Row().Scan(&stddev)
//
//	return stddev, err
//}
