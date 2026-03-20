package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) ListAPIKeys(ctx context.Context, params APIKeysListParams) ([]ManagedAPIKey, error) {
	query := url.Values{}
	if params.IncludeDisabled {
		query.Set("include_disabled", "true")
	}
	if params.Offset > 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	var res APIKeysResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/keys", query, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/keys", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) GetAPIKey(ctx context.Context, hash string) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) UpdateAPIKey(ctx context.Context, hash string, req UpdateAPIKeyRequest) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodPatch, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, hash string) error {
	return c.doJSON(ctx, http.MethodDelete, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, nil, nil)
}
