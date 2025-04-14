package entity_test

import (
	"api/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMLScore_NormalizeScore(t *testing.T) {
	mlScore := entity.MLScore{Score: 80}
	result := mlScore.NormalizeScore(decimal.NewFromFloat(70.0), decimal.NewFromFloat(10.0))
	assert.NotZero(t, result)
}

func TestMLScore_ToResponse(t *testing.T) {
	mlScore := entity.MLScore{
		ClientID:     "client-id",
		AdvertiserID: "advertiser-id",
		Score:        90,
		CreatedAt:    time.Now(),
	}

	response := mlScore.ToResponse()
	assert.Equal(t, "client-id", response.ClientID)
	assert.Equal(t, "advertiser-id", response.AdvertiserID)
	assert.Equal(t, 90, response.Score)
}
