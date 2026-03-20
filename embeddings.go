package gopenrouter

import (
	"context"
	"net/http"
)

func (c *Client) CreateEmbeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	var res EmbeddingResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/embeddings", nil, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ListEmbeddingsModels(ctx context.Context) (*ModelsList, error) {
	var res ModelsList
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/embeddings/models", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
