package api

import (
	"bot/pkg/advertising/errorz"
	"bot/pkg/advertising/response"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetClientByID retrieves a client by his ID.
func (c *Client) GetClientByID(clientID string) (response.Client, error) {
	url := fmt.Sprintf(getClientByIDEndpoint, clientID)

	resp, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return response.Client{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return response.Client{}, errorz.ErrInvalidData
	case http.StatusNotFound:
		return response.Client{}, errorz.ErrClientNotFound
	case http.StatusInternalServerError:
		return response.Client{}, errorz.ErrUnexpected
	}

	var clientResponse response.Client
	if err = json.NewDecoder(resp.Body).Decode(&clientResponse); err != nil {
		return response.Client{}, err
	}

	return clientResponse, nil
}
