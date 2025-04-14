package entity

import (
	"api/internal/entity/response"

	"github.com/shopspring/decimal"
)

// Statistic represents ad statistics entity.
type Statistic struct {
	ImpressionsCount int `json:"impressions_count"`
	ClicksCount      int `json:"clicks_count"`
	Conversion       float64
	SpentImpressions decimal.Decimal `json:"spent_impressions"`
	SpentClicks      decimal.Decimal `json:"spent_clicks"`
	SpentTotal       decimal.Decimal
	Date             int `json:"date"`
}

// ToGetStatistic converts Statistic into a response-friendly statistic structure.
func (e *Statistic) ToGetStatistic() response.GetStatistic {
	return response.GetStatistic{
		ImpressionsCount: e.ImpressionsCount,
		ClicksCount:      e.ClicksCount,
		Conversion:       e.Conversion,
		SpentImpressions: e.SpentImpressions.InexactFloat64(),
		SpentClicks:      e.SpentClicks.InexactFloat64(),
		SpentTotal:       e.SpentTotal.InexactFloat64(),
	}
}

// ToDailyStatistic converts Statistic into a response-friendly statistic structure.
func (e *Statistic) ToDailyStatistic() response.DailyStatistic {
	return response.DailyStatistic{
		ImpressionsCount: e.ImpressionsCount,
		ClicksCount:      e.ClicksCount,
		Conversion:       e.Conversion,
		SpentImpressions: e.SpentImpressions.InexactFloat64(),
		SpentClicks:      e.SpentClicks.InexactFloat64(),
		SpentTotal:       e.SpentTotal.InexactFloat64(),
		Date:             e.Date,
	}
}
