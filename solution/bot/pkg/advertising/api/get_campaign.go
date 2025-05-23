package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCampaign retrieves campaign with provided advertiser and campaign id.
func (c *Client) GetCampaign(advertiserID, campaignID string) (response.Campaign, error) {
	url := fmt.Sprintf(getCampaignEndpoint, advertiserID, campaignID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.Campaign{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Campaign{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Campaign{}, errorz.ErrCampaignNotFound
	case http.StatusInternalServerError:
		return response.Campaign{}, errorz.ErrUnexpected
	}

	var getCampaignResponse response.Campaign
	if err = json.NewDecoder(resp.Body).Decode(&getCampaignResponse); err != nil {
		return response.Campaign{}, err
	}

	return getCampaignResponse, nil
}
