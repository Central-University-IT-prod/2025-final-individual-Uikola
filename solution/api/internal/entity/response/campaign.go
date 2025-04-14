package response

// Campaign represents the response structure for a campaign entity.
type Campaign struct {
	CampaignID        string    `json:"campaign_id"`
	AdvertiserID      string    `json:"advertiser_id"`
	ImpressionsLimit  int       `json:"impressions_limit"`
	ClicksLimit       int       `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int       `json:"start_date"`
	EndDate           int       `json:"end_date"`
	ImageURL          string    `json:"image_url,omitempty"`
	PassedModeration  bool      `json:"passed_moderation"`
	Targeting         Targeting `json:"targeting"`
}

// Targeting represents the targeting structure for a campaign response structure.
type Targeting struct {
	Gender   *string `json:"gender,omitempty"`
	AgeFrom  *int    `json:"age_from,omitempty"`
	AgeTo    *int    `json:"age_to,omitempty"`
	Location *string `json:"location,omitempty"`
}

// GetAd represents the response structure for GetAd endpoint.
type GetAd struct {
	AdID         string `json:"ad_id"`
	AdTitle      string `json:"ad_title"`
	AdText       string `json:"ad_text"`
	AdvertiserID string `json:"advertiser_id"`
	ImageURL     string `json:"image_url,omitempty"`
}

// UploadImage represents the response structure for UploadImage endpoint.
type UploadImage struct {
	ImageURL string `json:"image_url"`
}

// GenerateAdText represents the response structure for GenerateAdText endpoint.
type GenerateAdText struct {
	GeneratedText string `json:"generated_text"`
}
