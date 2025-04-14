package redis

import (
	"bot/pkg/advertising/response"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AdvertiserRepository struct {
	redis *redis.Client
}

func NewAdvertiserRepository(client *redis.Client) *AdvertiserRepository {
	return &AdvertiserRepository{
		redis: client,
	}
}

func (s *AdvertiserRepository) Get(userID int64) (response.Advertiser, error) {
	advertiserBytes, err := s.redis.Get(context.Background(), fmt.Sprintf("%d", userID)).Result()
	if err != nil {
		return response.Advertiser{}, nil
	}

	var advertiser response.Advertiser
	if err = json.Unmarshal([]byte(advertiserBytes), &advertiser); err != nil {
		return response.Advertiser{}, err
	}

	return advertiser, nil
}

func (s *AdvertiserRepository) Set(userID int64, advertiser response.Advertiser, expiration time.Duration) {
	advertiserBytes, _ := json.Marshal(advertiser)
	s.redis.Set(context.Background(), fmt.Sprintf("%d", userID), advertiserBytes, expiration)
}

func (s *AdvertiserRepository) Clear(userID int64) {
	s.redis.Del(context.Background(), fmt.Sprintf("%d", userID))
}
