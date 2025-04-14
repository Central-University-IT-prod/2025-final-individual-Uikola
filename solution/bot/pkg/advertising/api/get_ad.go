package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetAd(clientID string) (response.GetAd, error) {
	url := fmt.Sprintf(getAdEndpoint, clientID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.GetAd{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return response.GetAd{}, errorz.ErrCampaignNotFound
	case http.StatusInternalServerError:
		return response.GetAd{}, errorz.ErrUnexpected
	}

	var getAdResponse response.GetAd
	if err = json.NewDecoder(resp.Body).Decode(&getAdResponse); err != nil {
		return response.GetAd{}, err
	}

	return getAdResponse, nil
}
