package api

import (
	"bot/pkg/advertising/errorz"
	"fmt"
	"net/http"
)

// DeleteCampaign delete campaign with provided advertiser and campaign id.
func (c *Client) DeleteCampaign(advertiserID, campaignID string) error {
	url := fmt.Sprintf(deleteCampaignEndpoint, advertiserID, campaignID)

	resp, err := c.sendRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return errorz.ErrInvalidData
	case http.StatusInternalServerError:
		return errorz.ErrUnexpected
	}

	return nil
}
