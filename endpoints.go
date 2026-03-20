package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) ListZDREndpoints(ctx context.Context) (*ZDREndpointsList, error) {
	var res ZDREndpointsList
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/endpoints/zdr", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ListModelEndpoints(ctx context.Context, author, slug string) (*ModelEndpoints, error) {
	var res ModelEndpointsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/models/"+url.PathEscape(author)+"/"+url.PathEscape(slug)+"/endpoints", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}
