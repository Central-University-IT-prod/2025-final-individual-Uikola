package request

import "github.com/go-playground/validator/v10"

type Advance struct {
	CurrentDate int `json:"current_date" validate:"required,gte=0"`
}

func (r *Advance) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
