package entity

import (
	"api/internal/entity/response"
	"time"
)

// Advertiser represents the advertiser entity.
type Advertiser struct {
	AdvertiserID string `gorm:"primaryKey;not null;type:uuid"`
	Name         string `gorm:"not null"`
	CreatedAt    time.Time
}

// ToResponse converts an Advertiser entity into a response-friendly structure.
func (e *Advertiser) ToResponse() response.Advertiser {
	return response.Advertiser{
		AdvertiserID: e.AdvertiserID,
		Name:         e.Name,
	}
}
