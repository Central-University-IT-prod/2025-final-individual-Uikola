package entity

import (
	"api/internal/entity/response"
	"api/internal/utils"
	"time"

	"github.com/shopspring/decimal"
)

// MLScore represents the ml-score entity.
type MLScore struct {
	ClientID     string `json:"client_id" gorm:"primaryKey;type:uuid"`
	AdvertiserID string `json:"advertiser_id" gorm:"primaryKey;type:uuid"`
	Score        int    `json:"score" gorm:"not null"`
	CreatedAt    time.Time
}

func (e *MLScore) NormalizeScore(avg, stddev decimal.Decimal) decimal.Decimal {
	return utils.SigmoidNormalization(decimal.NewFromInt(int64(e.Score)), avg, stddev)
}

// ToResponse converts a MLScore entity into a response-friendly structure.
func (e *MLScore) ToResponse() response.MLScore {
	return response.MLScore{
		ClientID:     e.ClientID,
		AdvertiserID: e.AdvertiserID,
		Score:        e.Score,
	}
}
