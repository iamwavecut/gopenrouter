package gopenrouter

import (
	"context"
	"net/http"
)

func (c *Client) GetCredits(ctx context.Context) (*Credits, error) {
	var res CreditsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/credits", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) CreateCoinbaseCharge(ctx context.Context, req CoinbaseChargeRequest) (*CoinbaseCharge, error) {
	var res CoinbaseChargeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/credits/coinbase", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}
