package response

// MLScore represents the response structure for ml-score entity.
type MLScore struct {
	ClientID     string `json:"client_id"`
	AdvertiserID string `json:"advertiser_id"`
	Score        int    `json:"score"`
}
