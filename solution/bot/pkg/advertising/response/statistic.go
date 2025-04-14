package response

// Statistic represents aggregated advertising statistics without date information.
type Statistic struct {
	ImpressionsCount int     `json:"impressions_count"`
	ClicksCount      int     `json:"clicks_count"`
	Conversion       float64 `json:"conversion"`
	SpentImpressions float64 `json:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total"`
}
