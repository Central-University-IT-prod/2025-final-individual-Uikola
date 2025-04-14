package e2e

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type CampaignCreate struct {
	ImpressionsLimit  int       `json:"impressions_limit"`
	ClicksLimit       int       `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int       `json:"start_date"`
	EndDate           int       `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

// CampaignUpdate — структура для обновления существующей кампании
type CampaignUpdate struct {
	ImpressionsLimit  int       `json:"impressions_limit"`
	ClicksLimit       int       `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int       `json:"start_date"`
	EndDate           int       `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

// Targeting — структура для таргетинга кампании
type Targeting struct {
	Gender   string `json:"gender"`
	AgeFrom  int    `json:"age_from"`
	AgeTo    int    `json:"age_to"`
	Location string `json:"location"`
}

func TestCampaignLifecycle(t *testing.T) {
	advertiserID := uuid.New().String()
	campaignID := ""

	t.Run("Create Advertiser", func(t *testing.T) {
		api.POST("/advertisers/bulk").
			WithJSON([]map[string]interface{}{
				{
					"advertiser_id": advertiserID,
					"name":          "Test Advertiser",
				},
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Create Campaign", func(t *testing.T) {
		campaign := CampaignCreate{
			ImpressionsLimit:  1000,
			ClicksLimit:       500,
			CostPerImpression: 0.5,
			CostPerClick:      1.0,
			AdTitle:           "Test Ad",
			AdText:            "This is a test ad",
			StartDate:         1,
			EndDate:           30,
			Targeting: Targeting{
				Gender:   "ALL",
				AgeFrom:  18,
				AgeTo:    65,
				Location: "Moscow",
			},
		}

		response := api.POST("/advertisers/{advertiserId}/campaigns", advertiserID).
			WithJSON(campaign).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		response.Value("campaign_id").NotNull()
		response.Value("advertiser_id").IsEqual(advertiserID)
		response.Value("impressions_limit").IsEqual(campaign.ImpressionsLimit)
		response.Value("clicks_limit").IsEqual(campaign.ClicksLimit)
		response.Value("cost_per_impression").IsEqual(campaign.CostPerImpression)
		response.Value("cost_per_click").IsEqual(campaign.CostPerClick)
		response.Value("ad_title").IsEqual(campaign.AdTitle)
		response.Value("ad_text").IsEqual(campaign.AdText)
		response.Value("ad_text").IsEqual(campaign.AdText)
		response.Value("start_date").IsEqual(campaign.StartDate)
		response.Value("end_date").IsEqual(campaign.EndDate)
		response.Value("targeting").NotNull()
		response.Value("targeting").Object().Value("gender").IsEqual(campaign.Targeting.Gender)
		response.Value("targeting").Object().Value("age_from").IsEqual(campaign.Targeting.AgeFrom)
		response.Value("targeting").Object().Value("age_to").IsEqual(campaign.Targeting.AgeTo)
		response.Value("targeting").Object().Value("location").IsEqual(campaign.Targeting.Location)

		campaignID = response.Value("campaign_id").String().Raw()
		assert.NotEmpty(t, campaignID, "Campaign ID не должен быть пустым")
	})

	t.Run("Get Campaign By ID", func(t *testing.T) {
		response := api.GET("/advertisers/{advertiserId}/campaigns/{campaignId}", advertiserID, campaignID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		response.Value("campaign_id").IsEqual(campaignID)
		response.Value("advertiser_id").IsEqual(advertiserID)
		response.Value("impressions_limit").IsEqual(1000)
		response.Value("clicks_limit").IsEqual(500)
		response.Value("cost_per_impression").IsEqual(0.5)
		response.Value("cost_per_click").IsEqual(1.0)
		response.Value("ad_title").IsEqual("Test Ad")
		response.Value("ad_text").IsEqual("This is a test ad")
		response.Value("start_date").IsEqual(1)
		response.Value("end_date").IsEqual(30)
		response.Value("targeting").NotNull()
		response.Value("targeting").Object().Value("gender").IsEqual("ALL")
		response.Value("targeting").Object().Value("age_from").IsEqual(18)
		response.Value("targeting").Object().Value("age_to").IsEqual(65)
		response.Value("targeting").Object().Value("location").IsEqual("Moscow")
	})

	t.Run("List Campaigns", func(t *testing.T) {
		response := api.GET("/advertisers/{advertiserId}/campaigns", advertiserID).
			WithQuery("size", 10).
			WithQuery("page", 1).
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		assert.Equal(t, len(response.Raw()), 1, "Список кампаний должен содержать одну кампанию")
	})

	t.Run("Update Campaign", func(t *testing.T) {
		campaignUpdate := CampaignUpdate{
			ImpressionsLimit:  1500,
			ClicksLimit:       600,
			CostPerImpression: 0.4,
			CostPerClick:      0.9,
			AdTitle:           "Updated Ad",
			AdText:            "This is an updated ad",
			StartDate:         1,
			EndDate:           31,
			Targeting: Targeting{
				Gender:   "FEMALE",
				AgeFrom:  25,
				AgeTo:    45,
				Location: "Saint Petersburg",
			},
		}

		response := api.PUT("/advertisers/{advertiserId}/campaigns/{campaignId}", advertiserID, campaignID).
			WithJSON(campaignUpdate).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		response.Value("campaign_id").IsEqual(campaignID)
		response.Value("advertiser_id").IsEqual(advertiserID)
		response.Value("impressions_limit").IsEqual(campaignUpdate.ImpressionsLimit)
		response.Value("clicks_limit").IsEqual(campaignUpdate.ClicksLimit)
		response.Value("cost_per_impression").IsEqual(campaignUpdate.CostPerImpression)
		response.Value("cost_per_click").IsEqual(campaignUpdate.CostPerClick)
		response.Value("ad_title").IsEqual(campaignUpdate.AdTitle)
		response.Value("ad_text").IsEqual(campaignUpdate.AdText)
		response.Value("ad_text").IsEqual(campaignUpdate.AdText)
		response.Value("start_date").IsEqual(campaignUpdate.StartDate)
		response.Value("end_date").IsEqual(campaignUpdate.EndDate)
		response.Value("targeting").NotNull()
		response.Value("targeting").Object().Value("gender").IsEqual(campaignUpdate.Targeting.Gender)
		response.Value("targeting").Object().Value("age_from").IsEqual(campaignUpdate.Targeting.AgeFrom)
		response.Value("targeting").Object().Value("age_to").IsEqual(campaignUpdate.Targeting.AgeTo)
		response.Value("targeting").Object().Value("location").IsEqual(campaignUpdate.Targeting.Location)
	})

	t.Run("Delete Campaign", func(t *testing.T) {
		api.DELETE("/advertisers/{advertiserId}/campaigns/{campaignId}", advertiserID, campaignID).
			Expect().
			Status(http.StatusNoContent)

		api.GET("/advertisers/{advertiserId}/campaigns/{campaignId}", advertiserID, campaignID).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Create Campaign For Non-Existing Advertiser", func(t *testing.T) {
		api.POST("/advertisers/{advertiserId}/campaigns", uuid.New().String()).
			WithJSON(CampaignCreate{
				ImpressionsLimit:  1000,
				ClicksLimit:       500,
				CostPerImpression: 0.5,
				CostPerClick:      1.0,
				AdTitle:           "Test Ad",
				AdText:            "This is a test ad",
				StartDate:         1,
				EndDate:           30,
				Targeting: Targeting{
					Gender:   "ALL",
					AgeFrom:  18,
					AgeTo:    65,
					Location: "Moscow",
				},
			}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Create Campaign With Invalid Data", func(t *testing.T) {
		api.POST("/advertisers/{advertiserId}/campaigns", advertiserID).
			WithJSON(CampaignCreate{
				ImpressionsLimit:  -1000,
				ClicksLimit:       500,
				CostPerImpression: 0.5,
				CostPerClick:      1.0,
				AdTitle:           "Invalid ad",
				AdText:            "",
				StartDate:         1,
				EndDate:           -10,
				Targeting: Targeting{
					Gender:   "TEST",
					AgeFrom:  18,
					AgeTo:    65,
					Location: "Moscow",
				},
			}).
			Expect().
			Status(http.StatusBadRequest)
	})
}
