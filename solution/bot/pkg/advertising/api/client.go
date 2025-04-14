package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	getCurrentDateEndpoint     = "time/advance"
	getClientByIDEndpoint      = "clients/%s"                               // clients/{{ clientID }}
	getAdvertiserByIDEndpoint  = "advertisers/%s"                           // advertisers/{{ advertiserID }}
	generateAdTextEndpoint     = "advertisers/%s/generate-ad-text"          // advertisers/{{ advertiserID }}/generate-ad-text
	createCampaignEndpoint     = "advertisers/%s/campaigns"                 // advertisers/{{ advertiserID }}/campaigns
	getCampaignEndpoint        = "advertisers/%s/campaigns/%s"              // advertisers/{{ advertiserID }}/campaigns/{{ campaignID }}
	listCampaignsEndpoint      = "advertisers/%s/campaigns?size=%d&page=%d" // advertisers/{{ advertiserID }}/campaigns
	updateCampaignEndpoint     = "advertisers/%s/campaigns/%s"              // advertisers/{{ advertiserID }}/campaigns/{{ campaignID }}
	deleteCampaignEndpoint     = "advertisers/%s/campaigns/%s"              // advertisers/{{ advertiserID }}/campaigns/{{ campaignID }}
	getCampaignStatsEndpoint   = "stats/campaigns/%s"                       // stats/campaigns/{{ campaignID }}
	getAdvertiserStatsEndpoint = "stats/advertisers/%s/campaigns"           // stats/advertisers/{{ advertiserID }}/campaigns
	getAdEndpoint              = "ads?client_id=%s"
	clickAdEndpoint            = "ads/%s/click" // ads/{{ campaign_id }}/click
)

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url:    url,
		client: &http.Client{},
	}
}

// sendRequest is a generic method to handle HTTP requests.
func (c *Client) sendRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.url, endpoint)

	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}
