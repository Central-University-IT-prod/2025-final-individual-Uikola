package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCampaignStats retrieves a campaign stats by his ID.
func (c *Client) GetCampaignStats(campaignID string) (response.Statistic, error) {
	url := fmt.Sprintf(getCampaignStatsEndpoint, campaignID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.Statistic{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Statistic{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Statistic{}, errorz.ErrCampaignNotFound
	case http.StatusInternalServerError:
		return response.Statistic{}, errorz.ErrUnexpected
	}

	var getCampaignStatsResponse response.Statistic
	if err = json.NewDecoder(resp.Body).Decode(&getCampaignStatsResponse); err != nil {
		return response.Statistic{}, err
	}

	return getCampaignStatsResponse, nil
}
