package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) GetUserActivity(ctx context.Context, params ActivityParams) ([]ActivityItem, error) {
	query := url.Values{}
	if params.Date != "" {
		query.Set("date", params.Date)
	}
	var res ActivityResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/activity", query, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
