package response

// GenerateAdText represents the GenerateAdText response structure.
type GenerateAdText struct {
	GeneratedText string `json:"generated_text"`
}

// Campaign represents the response structure for campaign.
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

// IsOver checks if the campaign is over on the given day.
func (r *Campaign) IsOver(currentDay int) bool {
	return r.EndDate < currentDay
}

// IsActive checks if the campaign is active on the given day.
func (r *Campaign) IsActive(currentDay int) bool {
	return (currentDay >= r.StartDate) && (currentDay <= r.EndDate)
}

// Targeting represents the response structure for targeting.
type Targeting struct {
	Gender   *string `json:"gender,omitempty"`
	AgeFrom  *int    `json:"age_from,omitempty"`
	AgeTo    *int    `json:"age_to,omitempty"`
	Location *string `json:"location,omitempty"`
}

type GetAd struct {
	AdID         string `json:"ad_id"`
	AdTitle      string `json:"ad_title"`
	AdText       string `json:"ad_text"`
	AdvertiserID string `json:"advertiser_id"`
	ImageURL     string `json:"image_url,omitempty"`
}
