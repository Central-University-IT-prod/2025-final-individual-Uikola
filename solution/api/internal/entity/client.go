package entity

import (
	"api/internal/entity/response"
	"time"
)

// Client represents the client entity.
type Client struct {
	ClientID  string       `gorm:"primaryKey;not null;type:uuid"`
	Login     string       `gorm:"not null"`
	Age       int          `gorm:"not null"`
	Location  string       `gorm:"not null"`
	Gender    ClientGender `gorm:"not null"`
	CreatedAt time.Time
}

// ToResponse converts a Client entity into a response-friendly structure.
func (e *Client) ToResponse() response.Client {
	return response.Client{
		ClientID: e.ClientID,
		Login:    e.Login,
		Age:      e.Age,
		Location: e.Location,
		Gender:   string(e.Gender),
	}
}

// ClientGender represents the client's gender as an enum-like string type.
type ClientGender string

const (
	ClientGenderMale   ClientGender = "MALE"
	ClientGenderFemale ClientGender = "FEMALE"
)
