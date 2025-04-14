package request

// GenerateAdText represents the GenerateAdText request structure.
type GenerateAdText struct {
	AdTitle string `json:"ad_title"`
	Context string `json:"context"`
}

// CreateCampaign represents the CreateCampaign request structure.
type CreateCampaign struct {
	ImpressionsLimit  int       `json:"impressions_limit"`
	ClicksLimit       int       `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int       `json:"start_date"`
	EndDate           int       `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

// UpdateCampaign represents the UpdateCampaign request structure.
type UpdateCampaign struct {
	ImpressionsLimit  int        `json:"impressions_limit"`
	ClicksLimit       int        `json:"clicks_limit"`
	CostPerImpression float64    `json:"cost_per_impression"`
	CostPerClick      float64    `json:"cost_per_click"`
	AdTitle           *string    `json:"ad_title"`
	AdText            *string    `json:"ad_text"`
	StartDate         int        `json:"start_date"`
	EndDate           int        `json:"end_date"`
	Targeting         *Targeting `json:"targeting"`
}

// Targeting represents the request structure for targeting.
type Targeting struct {
	Gender   *string `json:"gender"`
	AgeFrom  *int    `json:"age_from"`
	AgeTo    *int    `json:"age_to"`
	Location *string `json:"location"`
}

// ClickAd represents request body for click on add.
type ClickAd struct {
	ClientID string `json:"client_id"`
}
