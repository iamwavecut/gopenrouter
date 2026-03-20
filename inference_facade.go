package gopenrouter

import (
	"context"
	"net/http"

	anthropicpkg "github.com/iamwavecut/gopenrouter/anthropic"
	responsespkg "github.com/iamwavecut/gopenrouter/responses"
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

func (c *Client) CreateAnthropicMessage(ctx context.Context, req AnthropicMessageRequest) (*AnthropicMessageResponse, error) {
	if req.Stream {
		return nil, &RequestError{Err: errStreamMethodRequired("CreateAnthropicMessageStream")}
	}
	var res AnthropicMessageResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/messages", nil, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateAnthropicMessageStream(ctx context.Context, req AnthropicMessageRequest) (*AnthropicMessageStream, error) {
	req.Stream = true
	httpReq, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+"/messages", req)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.config.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Err: err}
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, decodeErrorResponse(resp)
	}
	return anthropicpkg.NewStream(resp), nil
}

func (c *Client) CreateResponse(ctx context.Context, req ResponseRequest) (*Response, error) {
	if req.Stream {
		return nil, &RequestError{Err: errStreamMethodRequired("CreateResponseStream")}
	}
	var res Response
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/responses", nil, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateResponseStream(ctx context.Context, req ResponseRequest) (*ResponseStream, error) {
	req.Stream = true
	httpReq, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+"/responses", req)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.config.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Err: err}
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, decodeErrorResponse(resp)
	}
	return responsespkg.NewStream(resp), nil
}

func errStreamMethodRequired(name string) error {
	return &RequestError{Err: &APIError{Message: "use " + name + " for streaming requests"}}
}
