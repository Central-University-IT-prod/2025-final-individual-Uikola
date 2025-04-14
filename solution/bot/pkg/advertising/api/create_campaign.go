package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateCampaign create campaign using provided request data.
func (c *Client) CreateCampaign(requestBody request.CreateCampaign, advertiserID string) (response.Campaign, error) {
	url := fmt.Sprintf(createCampaignEndpoint, advertiserID)

	resp, err := c.sendRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return response.Campaign{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Campaign{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Campaign{}, errorz.ErrAdvertiserNotFound
	case http.StatusInternalServerError:
		return response.Campaign{}, errorz.ErrUnexpected
	}

	var createCampaignResponse response.Campaign
	if err = json.NewDecoder(resp.Body).Decode(&createCampaignResponse); err != nil {
		return response.Campaign{}, err
	}

	return createCampaignResponse, nil
}
