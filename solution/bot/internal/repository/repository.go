package repository

import (
	"bot/internal/entity"
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
	"context"
	"time"
)

type UserRepository interface {
	Save(ctx context.Context, user entity.User) (entity.User, error)

	GetByID(ctx context.Context, userID int64) (entity.User, error)

	Delete(ctx context.Context, userID int64) error
}

type StateRepository interface {
	Get(userID int64) (string, error)

	Set(userID int64, state string, expiration time.Duration) error

	Delete(userID int64)
}

type CampaignRepository interface {
	Get(userID int64) (request.CreateCampaign, error)

	Set(userID int64, campaign request.CreateCampaign, expiration time.Duration)

	Clear(userID int64)
}

type AdvertiserRepository interface {
	Get(userID int64) (response.Advertiser, error)

	Set(userID int64, advertiser response.Advertiser, expiration time.Duration)

	Clear(userID int64)
}
