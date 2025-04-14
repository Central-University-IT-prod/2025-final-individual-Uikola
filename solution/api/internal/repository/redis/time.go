package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// DayKey is the Redis key used to store the day value.
const DayKey = "day"

// TimeRepository provides methods for interacting with Redis to store and retrieve time-related data.
type TimeRepository struct {
	redis *redis.Client
}

// NewTimeRepository creates a new instance of TimeRepository.
func NewTimeRepository(redis *redis.Client) *TimeRepository {
	return &TimeRepository{
		redis: redis,
	}
}

// Get retrieves the stored day value from Redis.
func (s *TimeRepository) Get(ctx context.Context) (int, error) {
	dayStr, err := s.redis.Get(ctx, DayKey).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	day, _ := strconv.Atoi(dayStr)
	return day, nil
}

// Set stores a day value in Redis with an expiration time.
func (s *TimeRepository) Set(ctx context.Context, day int, expiration time.Duration) {
	s.redis.Set(ctx, DayKey, day, expiration)
}
