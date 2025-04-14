package request

import (
	"api/internal/entity"

	"github.com/go-playground/validator/v10"
)

// MLScore represents the ml-score request body.
type MLScore struct {
	ClientID     string `json:"client_id" validate:"required,uuid"`
	AdvertiserID string `json:"advertiser_id" validate:"required,uuid"`
	Score        int    `json:"score" validate:"required,gte=0"`
}

// Validate validate ml-score request body.
func (r *MLScore) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// ToMLScore converts the request DTO (Data Transfer Object) into an entity.MLScore struct.
func (r *MLScore) ToMLScore() entity.MLScore {
	return entity.MLScore{
		ClientID:     r.ClientID,
		AdvertiserID: r.AdvertiserID,
		Score:        r.Score,
	}
}
