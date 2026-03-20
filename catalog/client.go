package catalog

import "context"

type Client struct {
	backend backend
}

type backend interface {
	ListModels(ctx context.Context) (*ModelsList, error)
	ListModelsWithParams(ctx context.Context, params ModelsListParams) (*ModelsList, error)
	CountModels(ctx context.Context) (int, error)
	ListModelsForUser(ctx context.Context) (*ModelsList, error)
	ListProviders(ctx context.Context) (*ProvidersList, error)
	ListZDREndpoints(ctx context.Context) (*ZDREndpointsList, error)
	ListModelEndpoints(ctx context.Context, author, slug string) (*ModelEndpoints, error)
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) ListModels(ctx context.Context) (*ModelsList, error) {
	return c.backend.ListModels(ctx)
}

func (c *Client) ListModelsWithParams(ctx context.Context, params ModelsListParams) (*ModelsList, error) {
	return c.backend.ListModelsWithParams(ctx, params)
}

func (c *Client) CountModels(ctx context.Context) (int, error) {
	return c.backend.CountModels(ctx)
}

func (c *Client) ListModelsForUser(ctx context.Context) (*ModelsList, error) {
	return c.backend.ListModelsForUser(ctx)
}

func (c *Client) ListProviders(ctx context.Context) (*ProvidersList, error) {
	return c.backend.ListProviders(ctx)
}

func (c *Client) ListZDREndpoints(ctx context.Context) (*ZDREndpointsList, error) {
	return c.backend.ListZDREndpoints(ctx)
}

func (c *Client) ListModelEndpoints(ctx context.Context, author, slug string) (*ModelEndpoints, error) {
	return c.backend.ListModelEndpoints(ctx, author, slug)
}
