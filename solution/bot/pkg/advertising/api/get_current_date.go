package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"net/http"
)

// GetCurrentDate retrieves a current api date.
func (c *Client) GetCurrentDate() (int, error) {
	url := getCurrentDateEndpoint

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return 0, errorz.ErrUnexpected
	}

	var getCurrentDateResponse response.GetCurrentDate
	if err = json.NewDecoder(resp.Body).Decode(&getCurrentDateResponse); err != nil {
		return 0, err
	}

	return getCurrentDateResponse.CurrentDate, nil
}
