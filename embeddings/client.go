package embeddings

import "context"

type Client struct {
	backend backend
}

type backend interface {
	CreateEmbeddings(ctx context.Context, req Request) (*Response, error)
	ListEmbeddingsModels(ctx context.Context) (*ModelsList, error)
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) Create(ctx context.Context, req Request) (*Response, error) {
	return c.backend.CreateEmbeddings(ctx, req)
}

func (c *Client) ListModels(ctx context.Context) (*ModelsList, error) {
	return c.backend.ListEmbeddingsModels(ctx)
}
