package oauth

import "context"

type Client struct {
	backend backend
}

type backend interface {
	CreateAuthCode(ctx context.Context, req CreateAuthCodeRequest) (*AuthCode, error)
	ExchangeAuthCodeForAPIKey(ctx context.Context, req ExchangeAuthCodeRequest) (*ExchangeAuthCodeResponse, error)
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) CreateCode(ctx context.Context, req CreateAuthCodeRequest) (*AuthCode, error) {
	return c.backend.CreateAuthCode(ctx, req)
}

func (c *Client) ExchangeCode(ctx context.Context, req ExchangeAuthCodeRequest) (*ExchangeAuthCodeResponse, error) {
	return c.backend.ExchangeAuthCodeForAPIKey(ctx, req)
}
