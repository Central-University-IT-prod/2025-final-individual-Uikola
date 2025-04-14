package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"fmt"
	"net/http"
)

func (c *Client) ClickAd(requestBody request.ClickAd, campaignID string) error {
	url := fmt.Sprintf(clickAdEndpoint, campaignID)

	resp, err := c.sendRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return errorz.ErrInvalidData
	case http.StatusNotFound:
		return errorz.ErrCampaignNotFound
	case http.StatusForbidden:
		return errorz.ErrImpressionNotFound
	case http.StatusInternalServerError:
		return errorz.ErrUnexpected
	}

	return nil
}
