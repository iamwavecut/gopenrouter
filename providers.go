package gopenrouter

import (
	"context"
	"net/http"
)

func (c *Client) ListProviders(ctx context.Context) (*ProvidersList, error) {
	var res ProvidersList
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/providers", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
