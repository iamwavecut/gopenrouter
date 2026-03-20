package gopenrouter

import (
	"context"
	"net/http"
)

func (c *Client) GetCurrentKey(ctx context.Context) (*KeyData, error) {
	var res KeyCheckResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/key", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Deprecated: use GetCurrentKey. This method is kept as a compatibility alias.
func (c *Client) CheckCredits(ctx context.Context) (*KeyData, error) {
	keyData, err := c.GetCurrentKey(ctx)
	if err == nil {
		return keyData, nil
	}

	var res KeyCheckResponse
	if fallbackErr := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/auth/key", nil, nil, &res); fallbackErr != nil {
		return nil, err
	}
	return &res.Data, nil
}
