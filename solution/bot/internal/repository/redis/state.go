package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateRepository struct {
	redis *redis.Client
}

func NewStateRepository(redis *redis.Client) *StateRepository {
	return &StateRepository{
		redis: redis,
	}
}

func (s *StateRepository) Get(userID int64) (string, error) {
	data, err := s.redis.Get(context.Background(), fmt.Sprintf("%d", userID)).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (s *StateRepository) Set(userID int64, state string, expiration time.Duration) error {
	return s.redis.Set(context.Background(), fmt.Sprintf("%d", userID), state, expiration).Err()
}

func (s *StateRepository) Delete(userID int64) {
	s.redis.Del(context.Background(), fmt.Sprintf("%d", userID))
}
