package henrik

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

func (c *Client) GetMatches(ctx context.Context, region, name, tag string) ([]models.Match, error) {
	endpoint := fmt.Sprintf("/valorant/v3/matches/%s/%s/%s?mode=competitive&size=10", region, name, tag)

	apiResp, err := c.callAPI(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var matches []models.Match
	if len(apiResp.Data) == 0 || string(apiResp.Data) == "null" {
		return []models.Match{}, nil
	}

	if err := json.Unmarshal(apiResp.Data, &matches); err != nil {
		return nil, fmt.Errorf("error parsing matches: %w", err)
	}

	return matches, nil
}
