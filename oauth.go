package gopenrouter

import (
	"context"
	"net/http"
)

func (c *Client) CreateAuthCode(ctx context.Context, req CreateAuthCodeRequest) (*AuthCode, error) {
	var res CreateAuthCodeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/auth/keys/code", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ExchangeAuthCodeForAPIKey(ctx context.Context, req ExchangeAuthCodeRequest) (*ExchangeAuthCodeResponse, error) {
	var res ExchangeAuthCodeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/auth/keys", nil, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
