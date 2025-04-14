package entity_test

import (
	"api/internal/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCampaign_IsActive(t *testing.T) {
	campaign := entity.Campaign{StartDate: 10, EndDate: 20}
	assert.True(t, campaign.IsActive(15))
	assert.False(t, campaign.IsActive(9))
	assert.False(t, campaign.IsActive(21))
}

func TestCampaign_IsOver(t *testing.T) {
	campaign := entity.Campaign{EndDate: 20}
	assert.True(t, campaign.IsOver(21))
	assert.False(t, campaign.IsOver(20))
	assert.False(t, campaign.IsOver(19))
}

func TestCampaign_NormalizeCPI(t *testing.T) {
	campaign := entity.Campaign{CostPerImpression: decimal.NewFromFloat(1.0)}
	result := campaign.NormalizeCPI(decimal.NewFromFloat(1.0), decimal.NewFromFloat(0.5))
	assert.NotZero(t, result)
}

func TestCampaign_NormalizeCPC(t *testing.T) {
	campaign := entity.Campaign{CostPerClick: decimal.NewFromFloat(1.0)}
	result := campaign.NormalizeCPC(decimal.NewFromFloat(1.0), decimal.NewFromFloat(0.5))
	assert.NotZero(t, result)
}

func TestCampaign_ToResponse(t *testing.T) {
	campaign := entity.Campaign{
		CampaignID:        "test-id",
		AdvertiserID:      "advertiser-id",
		ImpressionsLimit:  1000,
		ClicksLimit:       500,
		CostPerImpression: decimal.NewFromFloat(0.5),
		CostPerClick:      decimal.NewFromFloat(1.0),
		AdTitle:           "Test Ad",
		AdText:            "This is a test advertisement.",
		StartDate:         1,
		EndDate:           10,
		PassedModeration:  true,
	}

	response := campaign.ToResponse("http://example.com/image.png")
	assert.Equal(t, "test-id", response.CampaignID)
	assert.Equal(t, "advertiser-id", response.AdvertiserID)
	assert.Equal(t, 1000, response.ImpressionsLimit)
	assert.Equal(t, 500, response.ClicksLimit)
	assert.Equal(t, 0.5, response.CostPerImpression)
	assert.Equal(t, 1.0, response.CostPerClick)
	assert.Equal(t, "Test Ad", response.AdTitle)
	assert.Equal(t, "This is a test advertisement.", response.AdText)
	assert.Equal(t, 1, response.StartDate)
	assert.Equal(t, 10, response.EndDate)
	assert.Equal(t, "http://example.com/image.png", response.ImageURL)
	assert.True(t, response.PassedModeration)
}

func TestTargeting_ToResponse(t *testing.T) {
	gender := entity.TargetingGenderMale
	ageFrom := 18
	ageTo := 35
	location := "New York"
	targeting := entity.Targeting{
		Gender:   &gender,
		AgeFrom:  &ageFrom,
		AgeTo:    &ageTo,
		Location: &location,
	}

	response := targeting.ToResponse()
	assert.Equal(t, "MALE", *response.Gender)
	assert.Equal(t, 18, *response.AgeFrom)
	assert.Equal(t, 35, *response.AgeTo)
	assert.Equal(t, "New York", *response.Location)
}

func TestTargetingGender_ToString(t *testing.T) {
	gender := entity.TargetingGenderMale
	assert.Equal(t, "MALE", *gender.ToString())

	nilGender := (*entity.TargetingGender)(nil)
	assert.Nil(t, nilGender.ToString())
}
