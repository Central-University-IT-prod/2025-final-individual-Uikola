package request

import (
	"api/internal/entity"

	"github.com/go-playground/validator/v10"
)

// Advertiser represents the advertiser request body.
type Advertiser struct {
	AdvertiserID string `json:"advertiser_id" validate:"required,uuid"`
	Name         string `json:"name" validate:"required"`
}

// Validate validate advertiser request body.
func (r *Advertiser) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// ToAdvertiser converts the request DTO (Data Transfer Object) into an entity.Advertiser struct.
func (r *Advertiser) ToAdvertiser() entity.Advertiser {
	return entity.Advertiser{
		AdvertiserID: r.AdvertiserID,
		Name:         r.Name,
	}
}
