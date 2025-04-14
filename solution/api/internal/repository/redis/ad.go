package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	impressionsCountKey = "ad_impressions:%s"
	clicksCountKey      = "ad_clicks:%s"
	seenKeyPrefix       = "client_seen:%s"
	clickedKeyPrefix    = "client_clicked:%s"
)

type AdRepository struct {
	redis *redis.Client
}

func NewAdRepository(redis *redis.Client) *AdRepository {
	return &AdRepository{
		redis: redis,
	}
}

// TryIncrementImpressions attempts to increment the impressions count for a campaign.
// Returns true if incremented successfully (within the limit), otherwise false
func (r *AdRepository) TryIncrementImpressions(ctx context.Context, campaignID string, limit int) (bool, error) {
	key := fmt.Sprintf(impressionsCountKey, campaignID)
	maxAllowed := float64(limit) * 1.05

	luaScript := `
		local current = tonumber(redis.call('GET', KEYS[1]) or '0')
		if (current + 1) < tonumber(ARGV[1]) then
			return redis.call('INCR', KEYS[1])
		else
			return -1
		end
	`

	result, err := r.redis.Eval(ctx, luaScript, []string{key}, maxAllowed).Int()
	if err != nil {
		return false, err
	}

	return result > -1, nil
}

// IncrementClick increments the clicks count for a given campaign.
func (r *AdRepository) IncrementClick(ctx context.Context, campaignID string) error {
	key := fmt.Sprintf(clicksCountKey, campaignID)
	_, err := r.redis.Incr(ctx, key).Result()
	return err
}

// GetImpressionsCount retrieves the impressions count for a campaign.
// Returns 0 if the key is not found.
func (r *AdRepository) GetImpressionsCount(ctx context.Context, campaignID string) (int, error) {
	key := fmt.Sprintf(impressionsCountKey, campaignID)
	count, err := r.redis.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return count, err
}

// GetClicksCount retrieves the clicks count for a campaign.
// Returns 0 if the key is not found.
func (r *AdRepository) GetClicksCount(ctx context.Context, campaignID string) (int, error) {
	key := fmt.Sprintf(clicksCountKey, campaignID)
	count, err := r.redis.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return count, err
}

// AddSeenCampaign records that a client has seen a specific campaign.
func (r *AdRepository) AddSeenCampaign(ctx context.Context, clientID, campaignID string) error {
	key := fmt.Sprintf(seenKeyPrefix, clientID)
	return r.redis.SAdd(ctx, key, campaignID).Err()
}

// HasSeenCampaign checks if a client has already seen a specific campaign.
func (r *AdRepository) HasSeenCampaign(ctx context.Context, clientID, campaignID string) (bool, error) {
	key := fmt.Sprintf(seenKeyPrefix, clientID)
	return r.redis.SIsMember(ctx, key, campaignID).Result()
}

// AddClickedCampaign records that a client has clicked on a specific campaign.
func (r *AdRepository) AddClickedCampaign(ctx context.Context, clientID, campaignID string) error {
	key := fmt.Sprintf(clickedKeyPrefix, clientID)
	return r.redis.SAdd(ctx, key, campaignID).Err()
}

// HasClickedCampaign checks if a client has already clicked on a specific campaign.
func (r *AdRepository) HasClickedCampaign(ctx context.Context, clientID, campaignID string) (bool, error) {
	key := fmt.Sprintf(clickedKeyPrefix, clientID)
	return r.redis.SIsMember(ctx, key, campaignID).Result()
}
