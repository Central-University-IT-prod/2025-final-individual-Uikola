package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ListCampaigns retrieves a list of the advertiser campaigns.
func (c *Client) ListCampaigns(advertiserID string, size, page int) ([]response.Campaign, int, error) {
	url := fmt.Sprintf(listCampaignsEndpoint, advertiserID, size, page)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return []response.Campaign{}, 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return []response.Campaign{}, 0, errorz.ErrInvalidData
	case http.StatusNotFound:
		return []response.Campaign{}, 0, errorz.ErrAdvertiserNotFound
	case http.StatusInternalServerError:
		return []response.Campaign{}, 0, errorz.ErrUnexpected
	}

	var getCampaignResponse []response.Campaign
	if err = json.NewDecoder(resp.Body).Decode(&getCampaignResponse); err != nil {
		return []response.Campaign{}, 0, err
	}

	countStr := resp.Header.Get("X-Total-Count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return []response.Campaign{}, 0, err
	}

	return getCampaignResponse, count, nil
}
