package request

import (
	"api/internal/entity"
	"errors"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"

	"github.com/shopspring/decimal"
)

var errInvalidDate = errors.New("invalid date")

// CreateCampaign represents the request body for campaign creation.
type CreateCampaign struct {
	ImpressionsLimit  int       `json:"impressions_limit" validate:"required,gt=0,gtefield=ClicksLimit"`
	ClicksLimit       int       `json:"clicks_limit" validate:"required,gt=0"`
	CostPerImpression float64   `json:"cost_per_impression" validate:"required,gt=0"`
	CostPerClick      float64   `json:"cost_per_click" validate:"required,gt=0"`
	AdTitle           string    `json:"ad_title" validate:"required"`
	AdText            string    `json:"ad_text" validate:"required"`
	StartDate         int       `json:"start_date" validate:"required,gte=0"`
	EndDate           int       `json:"end_date" validate:"required,gte=0,gtfield=StartDate"`
	Targeting         Targeting `json:"targeting" validate:"omitempty"`
}

// Validate validate create campaign request body.
func (r *CreateCampaign) Validate(currentDay int) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	if !validateDate(currentDay, r.StartDate, r.EndDate) {
		return errInvalidDate
	}

	return nil
}

// ToCampaign converts the request DTO (Data Transfer Object) into an entity.Campaign struct.
func (r *CreateCampaign) ToCampaign(advertiserID string) entity.Campaign {
	return entity.Campaign{
		AdvertiserID:      advertiserID,
		ImpressionsLimit:  r.ImpressionsLimit,
		ClicksLimit:       r.ClicksLimit,
		CostPerImpression: decimal.NewFromFloat(r.CostPerImpression),
		CostPerClick:      decimal.NewFromFloat(r.CostPerClick),
		AdTitle:           r.AdTitle,
		AdText:            r.AdText,
		StartDate:         r.StartDate,
		EndDate:           r.EndDate,
		Targeting:         r.Targeting.ToTargeting(),
	}
}

// UpdateCampaign represents the request body for campaign update.
type UpdateCampaign struct {
	ImpressionsLimit  int        `json:"impressions_limit" validate:"required,gt=0,gtefield=ClicksLimit"`
	ClicksLimit       int        `json:"clicks_limit" validate:"required,gt=0"`
	CostPerImpression float64    `json:"cost_per_impression" validate:"required,gt=0"`
	CostPerClick      float64    `json:"cost_per_click" validate:"required,gt=0"`
	AdTitle           *string    `json:"ad_title" validate:"omitempty"`
	AdText            *string    `json:"ad_text" validate:"omitempty"`
	StartDate         int        `json:"start_date" validate:"required,gte=0"`
	EndDate           int        `json:"end_date" validate:"required,gte=0,gtfield=StartDate"`
	Targeting         *Targeting `json:"targeting" validate:"omitempty"`
}

// Validate validate create campaign request body.
func (r *UpdateCampaign) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// Targeting represents the targeting request body.
type Targeting struct {
	Gender   *entity.TargetingGender `json:"gender" validate:"omitempty,oneof=MALE FEMALE ALL"`
	AgeFrom  *int                    `json:"age_from" validate:"omitempty,gte=0"`
	AgeTo    *int                    `json:"age_to" validate:"omitempty,gte=0,gtefield=AgeFrom"`
	Location *string                 `json:"location" validate:"omitempty"`
}

// ToTargeting converts the request DTO (Data Transfer Object) into an entity.Targeting struct.
func (r *Targeting) ToTargeting() entity.Targeting {
	return entity.Targeting{
		TargetingID: uuid.New().String(),
		Gender:      r.Gender,
		AgeFrom:     r.AgeFrom,
		AgeTo:       r.AgeTo,
		Location:    r.Location,
	}
}

// ClickAd represents request body for click on add.
type ClickAd struct {
	ClientID string `json:"client_id" validate:"required"`
}

// Validate validate request body for click on ad.
func (r *ClickAd) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// GenerateAdText represents request body for generate ad text request.
type GenerateAdText struct {
	AdTitle string `json:"ad_title" validate:"required"`
	Context string `json:"context" validate:"omitempty"`
}

// Validate validate request body for generate ad text request.
func (r *GenerateAdText) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// Moderate represents request body for moderate campaign request.
type Moderate struct {
	PassedModeration *bool `json:"passed_moderation" validate:"required"`
}

// Validate validate request body for moderate campaign request.
func (r *Moderate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
