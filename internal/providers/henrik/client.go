package henrik

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	apiKey string
	client *resty.Client
}

func New(httpClient *http.Client, apiKey string) Client {
	client := resty.NewWithClient(httpClient).
		SetBaseURL("https://api.henrikdev.xyz")
	return Client{
		apiKey: apiKey,
		client: client,
	}
}

func (c *Client) callAPI(ctx context.Context, endpoint string) (*models.APIResponse, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Authorization", c.apiKey).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	if apiResp.Status != 200 && apiResp.Status != 0 {
		return nil, fmt.Errorf("API error: status %d", apiResp.Status)
	}

	return &apiResp, nil
}
