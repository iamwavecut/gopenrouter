package gopenrouter

import (
	"context"
	"net/http"

	anthropicpkg "github.com/iamwavecut/gopenrouter/anthropic"
)

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
