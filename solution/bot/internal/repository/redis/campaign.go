package redis

import (
	"bot/pkg/advertising/request"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CampaignRepository struct {
	redis *redis.Client
}

func NewCampaignRepository(client *redis.Client) *CampaignRepository {
	return &CampaignRepository{
		redis: client,
	}
}

func (s *CampaignRepository) Get(userID int64) (request.CreateCampaign, error) {
	campaignBytes, err := s.redis.Get(context.Background(), fmt.Sprintf("%d", userID)).Result()
	if err != nil {
		return request.CreateCampaign{}, nil
	}

	var event request.CreateCampaign
	if err = json.Unmarshal([]byte(campaignBytes), &event); err != nil {
		return request.CreateCampaign{}, err
	}

	return event, nil
}

func (s *CampaignRepository) Set(userID int64, campaign request.CreateCampaign, expiration time.Duration) {
	campaignBytes, _ := json.Marshal(campaign)
	s.redis.Set(context.Background(), fmt.Sprintf("%d", userID), campaignBytes, expiration)
}

func (s *CampaignRepository) Clear(userID int64) {
	s.redis.Del(context.Background(), fmt.Sprintf("%d", userID))
}
