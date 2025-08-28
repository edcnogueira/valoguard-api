package henrik

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

func (c *Client) GetMMR(ctx context.Context, region, name, tag string) (*models.MMR, error) {
	endpoint := fmt.Sprintf("/valorant/v2/mmr/%s/%s/%s", region, name, tag)

	apiResp, err := c.callAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var mmr models.MMR
	if len(apiResp.Data) == 0 || string(apiResp.Data) == "null" {
		mmr.CurrentTierPatched = "Unranked"
		return &mmr, nil
	}

	if err := json.Unmarshal(apiResp.Data, &mmr); err != nil {
		return nil, fmt.Errorf("error parsing MMR: %w", err)
	}

	return &mmr, nil
}
