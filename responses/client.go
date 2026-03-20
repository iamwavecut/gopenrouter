package responses

import "context"

type Client struct {
	backend backend
}

type backend interface {
	CreateResponse(ctx context.Context, req Request) (*Response, error)
	CreateResponseStream(ctx context.Context, req Request) (*Stream, error)
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) Create(ctx context.Context, req Request) (*Response, error) {
	return c.backend.CreateResponse(ctx, req)
}

func (c *Client) CreateStream(ctx context.Context, req Request) (*Stream, error) {
	return c.backend.CreateResponseStream(ctx, req)
}
