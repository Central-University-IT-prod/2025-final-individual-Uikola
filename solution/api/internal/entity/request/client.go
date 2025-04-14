package request

import (
	"api/internal/entity"

	"github.com/go-playground/validator/v10"
)

// Client represents the client request body.
type Client struct {
	ClientID string              `json:"client_id" validate:"required,uuid"`
	Login    string              `json:"login" validate:"required"`
	Age      int                 `json:"age" validate:"required"`
	Location string              `json:"location" validate:"required"`
	Gender   entity.ClientGender `json:"gender" validate:"required,oneof=MALE FEMALE"`
}

// Validate validate client request body.
func (r *Client) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}

// ToClient converts the request DTO (Data Transfer Object) into an entity.Client struct.
func (r *Client) ToClient() entity.Client {
	return entity.Client{
		ClientID: r.ClientID,
		Login:    r.Login,
		Age:      r.Age,
		Location: r.Location,
		Gender:   r.Gender,
	}
}
