package e2e

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type AdvertiserUpsert struct {
	AdvertiserID string `json:"advertiser_id"`
	Name         string `json:"name"`
}

func TestAdvertiserLifecycle(t *testing.T) {
	advertiserID := uuid.New().String()

	t.Run("Create Advertiser", func(t *testing.T) {
		advertisers := []AdvertiserUpsert{
			{
				AdvertiserID: advertiserID,
				Name:         "test_advertiser",
			},
		}

		response := api.POST("/advertisers/bulk").
			WithJSON(advertisers).
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		assert.Equal(t, 1, len(response.Raw()), "Ожидался ответ с одним рекламодателем")
		response.Value(0).Object().Value("advertiser_id").IsEqual(advertiserID)
	})

	t.Run("Get Advertiser By ID", func(t *testing.T) {
		advertiser := api.GET("/advertisers/{advertiserId}", advertiserID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		advertiser.Value("advertiser_id").IsEqual(advertiserID)
		advertiser.Value("name").IsEqual("test_advertiser")
	})

	t.Run("Update Advertiser", func(t *testing.T) {
		advertisers := []AdvertiserUpsert{
			{
				AdvertiserID: advertiserID,
				Name:         "updated_test_advertiser",
			},
		}

		api.POST("/advertisers/bulk").
			WithJSON(advertisers).
			Expect().
			Status(http.StatusOK)

		advertiser := api.GET("/advertisers/{advertiserId}", advertiserID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		advertiser.Value("name").IsEqual("updated_test_advertiser")
	})

	t.Run("Get Non-Existing Advertiser", func(t *testing.T) {
		api.GET("/advertisers/{advertiserId}", uuid.New().String()).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Invalid Advertiser Data", func(t *testing.T) {
		invalidAdvertiser := []map[string]interface{}{
			{
				"client_id": "",
				"name":      "invalid_advertiser",
			},
		}

		api.POST("/advertisers/bulk").
			WithJSON(invalidAdvertiser).
			Expect().
			Status(http.StatusBadRequest)
	})
}
