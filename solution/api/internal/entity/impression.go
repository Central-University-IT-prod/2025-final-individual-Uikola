package entity

import "github.com/shopspring/decimal"

// Impression represents an ad impression event in the ClickHouse database.
// It logs each time an ad is shown to a client, including the associated advertiser,
// campaign, client ID, and the cost per impression.
type Impression struct {
	AdvertiserID      string          `gorm:"primaryKey;type:UUID"`
	CampaignID        string          `gorm:"primaryKey;type:UUID"`
	ClientID          string          `gorm:"primaryKey;type:UUID"`
	CostPerImpression decimal.Decimal `gorm:"not null;type:DECIMAL(12,4)"`
	CreatedAt         int             `gorm:"not null"`
}

type ImpressionCount struct {
	CampaignID       string `gorm:"column:campaign_id"`
	ImpressionsCount int    `gorm:"column:impressions_count"`
}
