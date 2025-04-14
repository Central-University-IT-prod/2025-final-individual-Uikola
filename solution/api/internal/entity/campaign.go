package entity

import (
	"api/internal/entity/response"
	"api/internal/utils"
	"time"

	"gorm.io/gorm"

	"github.com/shopspring/decimal"
)

// Campaign represents the campaign entity.
type Campaign struct {
	CampaignID        string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AdvertiserID      string          `gorm:"not null;type:uuid"`
	ImpressionsLimit  int             `gorm:"not null"`
	ClicksLimit       int             `gorm:"not null"`
	CostPerImpression decimal.Decimal `gorm:"not null;type:decimal(12,4)"`
	CostPerClick      decimal.Decimal `gorm:"not null;type:decimal(12,4)"`
	AdTitle           string          `gorm:"not null"`
	AdText            string          `gorm:"not null"`
	StartDate         int             `gorm:"not null"`
	EndDate           int             `gorm:"not null"`
	Targeting         Targeting       `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE"`
	ImageID           string
	PassedModeration  bool `gorm:"default:false"`
	CreatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

// IsActive checks if the campaign is active on the given day.
func (e *Campaign) IsActive(currentDay int) bool {
	return (currentDay >= e.StartDate) && (currentDay <= e.EndDate)
}

func (e *Campaign) NormalizeCPI(avg, stddev decimal.Decimal) decimal.Decimal {
	return utils.SigmoidNormalization(e.CostPerImpression, avg, stddev)
}

func (e *Campaign) NormalizeCPC(avg, stddev decimal.Decimal) decimal.Decimal {
	return utils.SigmoidNormalization(e.CostPerClick, avg, stddev)
}

// IsOver checks if the campaign is over on the given day.
func (e *Campaign) IsOver(currentDay int) bool {
	return e.EndDate < currentDay
}

// ToResponse converts a Campaign entity into a response-friendly structure.
func (e *Campaign) ToResponse(imageURL string) response.Campaign {
	return response.Campaign{
		CampaignID:        e.CampaignID,
		AdvertiserID:      e.AdvertiserID,
		ImpressionsLimit:  e.ImpressionsLimit,
		ClicksLimit:       e.ClicksLimit,
		CostPerImpression: e.CostPerImpression.InexactFloat64(),
		CostPerClick:      e.CostPerClick.InexactFloat64(),
		AdTitle:           e.AdTitle,
		AdText:            e.AdText,
		StartDate:         e.StartDate,
		EndDate:           e.EndDate,
		ImageURL:          imageURL,
		PassedModeration:  e.PassedModeration,
		Targeting:         e.Targeting.ToResponse(),
	}
}

// Targeting represents the targeting entity.
type Targeting struct {
	TargetingID string `gorm:"primaryKey;type:uuid"`
	CampaignID  string `gorm:"not null;type:uuid"`
	Gender      *TargetingGender
	AgeFrom     *int
	AgeTo       *int
	Location    *string
	DeletedAd   gorm.DeletedAt
}

// TargetingGender represents the targeting gender as an enum-like string type.
type TargetingGender string

// ToString converts a TargetingGender pointer to a string pointer.
func (g *TargetingGender) ToString() *string {
	if g == nil {
		return nil
	}
	s := string(*g)
	return &s
}

const (
	TargetingGenderMale   TargetingGender = "MALE"
	TargetingGenderFemale TargetingGender = "FEMALE"
	TargetingGenderAll    TargetingGender = "ALL"
)

// ToResponse converts a Targeting entity into a response-friendly structure.
func (e *Targeting) ToResponse() response.Targeting {
	return response.Targeting{
		Gender:   e.Gender.ToString(),
		AgeFrom:  e.AgeFrom,
		AgeTo:    e.AgeTo,
		Location: e.Location,
	}
}
