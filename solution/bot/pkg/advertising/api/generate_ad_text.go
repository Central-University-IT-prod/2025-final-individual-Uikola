package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/request"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GenerateAdText generates an advertisement text using the provided request data.
func (c *Client) GenerateAdText(requestBody request.GenerateAdText, advertiserID string) (response.GenerateAdText, error) {
	url := fmt.Sprintf(generateAdTextEndpoint, advertiserID)

	resp, err := c.sendRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return response.GenerateAdText{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.GenerateAdText{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.GenerateAdText{}, errorz.ErrAdvertiserNotFound
	case http.StatusInternalServerError:
		return response.GenerateAdText{}, errorz.ErrUnexpected
	}

	var generateAdTextResponse response.GenerateAdText
	if err = json.NewDecoder(resp.Body).Decode(&generateAdTextResponse); err != nil {
		return response.GenerateAdText{}, err
	}

	return generateAdTextResponse, nil
}
