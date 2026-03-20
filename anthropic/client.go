package anthropic

import "context"

type Client struct {
	backend backend
}

type backend interface {
	CreateAnthropicMessage(ctx context.Context, req Request) (*Response, error)
	CreateAnthropicMessageStream(ctx context.Context, req Request) (*Stream, error)
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) Create(ctx context.Context, req Request) (*Response, error) {
	return c.backend.CreateAnthropicMessage(ctx, req)
}

func (c *Client) CreateStream(ctx context.Context, req Request) (*Stream, error) {
	return c.backend.CreateAnthropicMessageStream(ctx, req)
}
