package e2e

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
)

type MLScore struct {
	ClientID     string `json:"client_id" validate:"required,uuid"`
	AdvertiserID string `json:"advertiser_id" validate:"required,uuid"`
	Score        int    `json:"score" validate:"required,gte=0"`
}

func TestMLScoreLifecycle(t *testing.T) {
	clientID := uuid.New().String()
	advertiserID := uuid.New().String()

	t.Run("Create Client and Advertiser", func(t *testing.T) {
		api.POST("/clients/bulk").
			WithJSON([]ClientUpsert{
				{
					ClientID: clientID,
					Login:    "ml_client",
					Age:      28,
					Location: "New York",
					Gender:   "FEMALE",
				},
			}).
			Expect().
			Status(http.StatusOK)

		api.POST("/advertisers/bulk").
			WithJSON([]map[string]interface{}{
				{
					"advertiser_id": advertiserID,
					"name":          "ML Advertiser",
				},
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Create ML Score", func(t *testing.T) {
		score := MLScore{
			ClientID:     clientID,
			AdvertiserID: advertiserID,
			Score:        85,
		}

		api.POST("/ml-scores").
			WithJSON(score).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Update ML Score", func(t *testing.T) {
		score := MLScore{
			ClientID:     clientID,
			AdvertiserID: advertiserID,
			Score:        95,
		}

		api.POST("/ml-scores").
			WithJSON(score).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Invalid ML Score", func(t *testing.T) {
		score := MLScore{
			ClientID:     clientID,
			AdvertiserID: advertiserID,
			Score:        -10,
		}

		api.POST("/ml-scores").
			WithJSON(score).
			Expect().
			Status(http.StatusBadRequest)
	})
}
