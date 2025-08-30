package henrik

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

func (c *Client) GetAccount(ctx context.Context, name, tag string) (*models.Account, error) {
	endpoint := fmt.Sprintf("/valorant/v2/account/%s/%s", name, tag)
	
	apiResp, err := c.callAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	
	if len(apiResp.Data) == 0 || string(apiResp.Data) == "null" {
		return nil, fmt.Errorf("account data not found")
	}
	
	var account models.Account
	if err := json.Unmarshal(apiResp.Data, &account); err != nil {
		return nil, fmt.Errorf("error parsing account: %w", err)
	}
	
	return &account, nil
}