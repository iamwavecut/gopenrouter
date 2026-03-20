package gopenrouter

import (
	"context"
	"net/http"

	responsespkg "github.com/iamwavecut/gopenrouter/responses"
)

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
