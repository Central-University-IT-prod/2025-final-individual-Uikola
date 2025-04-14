package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateCampaign update campaign with provided advertiser and campaign id.
func (c *Client) UpdateCampaign(requestBody request.UpdateCampaign, advertiserID, campaignID string) (response.Campaign, error) {
	url := fmt.Sprintf(updateCampaignEndpoint, advertiserID, campaignID)

	resp, err := c.sendRequest(http.MethodPut, url, requestBody)
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

	var updateCampaignResponse response.Campaign
	if err = json.NewDecoder(resp.Body).Decode(&updateCampaignResponse); err != nil {
		return response.Campaign{}, err
	}

	return updateCampaignResponse, nil
}
