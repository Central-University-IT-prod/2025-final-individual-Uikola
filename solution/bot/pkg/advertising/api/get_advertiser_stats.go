package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAdvertiserStats retrieves an advertiser stats by his ID.
func (c *Client) GetAdvertiserStats(advertiserID string) (response.Statistic, error) {
	url := fmt.Sprintf(getAdvertiserStatsEndpoint, advertiserID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.Statistic{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Statistic{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Statistic{}, errorz.ErrAdvertiserNotFound
	case http.StatusInternalServerError:
		return response.Statistic{}, errorz.ErrUnexpected
	}

	var getAdvertiserStatsResponse response.Statistic
	if err = json.NewDecoder(resp.Body).Decode(&getAdvertiserStatsResponse); err != nil {
		return response.Statistic{}, err
	}

	return getAdvertiserStatsResponse, nil
}
