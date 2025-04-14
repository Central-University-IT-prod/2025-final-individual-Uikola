package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAdvertiserByID retrieves an advertiser by his ID.
func (c *Client) GetAdvertiserByID(advertiserID string) (response.Advertiser, error) {
	url := fmt.Sprintf(getAdvertiserByIDEndpoint, advertiserID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.Advertiser{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Advertiser{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Advertiser{}, errorz.ErrAdvertiserNotFound
	case http.StatusInternalServerError:
		return response.Advertiser{}, errorz.ErrUnexpected
	}

	var advertiserResponse response.Advertiser
	if err = json.NewDecoder(resp.Body).Decode(&advertiserResponse); err != nil {
		return response.Advertiser{}, err
	}

	return advertiserResponse, nil
}
