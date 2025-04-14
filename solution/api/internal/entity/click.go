package entity

import "github.com/shopspring/decimal"

// Click represents a recorded click event in the ClickHouse database.
// It stores information about a client clicking on an advertisement,
// including the associated advertiser, campaign, client ID, and cost per click.
type Click struct {
	AdvertiserID string          `gorm:"primaryKey;type:UUID"`
	CampaignID   string          `gorm:"primaryKey;type:UUID"`
	ClientID     string          `gorm:"primaryKey;type:UUID"`
	CostPerClick decimal.Decimal `gorm:"not null;type:DECIMAL(12,4)"`
	CreatedAt    int             `gorm:"not null"`
}

type ClickCount struct {
	CampaignID  string `gorm:"column:campaign_id"`
	ClicksCount int    `gorm:"column:clicks_count"`
}
