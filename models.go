package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) ListModels(ctx context.Context) (*ModelsList, error) {
	return c.ListModelsWithParams(ctx, ModelsListParams{})
}

func (c *Client) ListModelsWithParams(ctx context.Context, params ModelsListParams) (*ModelsList, error) {
	query := url.Values{}
	if params.Category != "" {
		query.Set("category", params.Category)
	}
	if params.SupportedParameters != "" {
		query.Set("supported_parameters", params.SupportedParameters)
	}

	var res ModelsList
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/models", query, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CountModels(ctx context.Context) (int, error) {
	var res ModelsCountResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/models/count", nil, nil, &res); err != nil {
		return 0, err
	}
	return res.Data.Count, nil
}

func (c *Client) ListModelsForUser(ctx context.Context) (*ModelsList, error) {
	var res ModelsList
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/models/user", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
