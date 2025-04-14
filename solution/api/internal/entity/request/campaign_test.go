package request_test

import (
	"api/internal/entity/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCampaign_Validate(t *testing.T) {
	validCampaign := request.CreateCampaign{
		ImpressionsLimit:  1000,
		ClicksLimit:       500,
		CostPerImpression: 0.5,
		CostPerClick:      1.0,
		AdTitle:           "Test Ad",
		AdText:            "This is a test advertisement.",
		StartDate:         1,
		EndDate:           10,
	}
	assert.NoError(t, validCampaign.Validate(0))

	invalidCampaign := request.CreateCampaign{AdTitle: "Invalid"}
	assert.Error(t, invalidCampaign.Validate(0))
}

func TestCreateCampaign_ToCampaign(t *testing.T) {
	r := request.CreateCampaign{
		ImpressionsLimit:  1000,
		ClicksLimit:       500,
		CostPerImpression: 0.5,
		CostPerClick:      1.0,
		AdTitle:           "Test Ad",
		AdText:            "This is a test advertisement.",
	}

	campaign := r.ToCampaign("advertiser-id")
	assert.Equal(t, "advertiser-id", campaign.AdvertiserID)
	assert.Equal(t, 1000, campaign.ImpressionsLimit)
}

func TestUpdateCampaign_Validate(t *testing.T) {
	r := request.UpdateCampaign{
		ImpressionsLimit:  1000,
		ClicksLimit:       500,
		CostPerImpression: 0.5,
		CostPerClick:      1.0,
		StartDate:         1,
		EndDate:           10,
	}
	assert.NoError(t, r.Validate())
}

func TestTargeting_ToTargeting(t *testing.T) {
	r := request.Targeting{Location: &[]string{"New York"}[0]}
	targeting := r.ToTargeting()
	assert.Equal(t, "New York", *targeting.Location)
}

func TestClickAd_Validate(t *testing.T) {
	r := request.ClickAd{ClientID: "client-id"}
	assert.NoError(t, r.Validate())

	invalid := request.ClickAd{}
	assert.Error(t, invalid.Validate())
}

func TestGenerateAdText_Validate(t *testing.T) {
	r := request.GenerateAdText{AdTitle: "Test Title"}
	assert.NoError(t, r.Validate())

	invalid := request.GenerateAdText{}
	assert.Error(t, invalid.Validate())
}

func TestModerate_Validate(t *testing.T) {
	passed := true
	r := request.Moderate{PassedModeration: &passed}
	assert.NoError(t, r.Validate())

	invalid := request.Moderate{}
	assert.Error(t, invalid.Validate())
}
