package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// KeyCheckResponse is the response from the /auth/key endpoint.
type KeyCheckResponse struct {
	Data KeyData `json:"data"`
}

// KeyData contains the details for an API key.
type KeyData struct {
	Label      string  `json:"label"`
	Usage      float64 `json:"usage"`
	Limit      float64 `json:"limit"`
	IsFreeTier bool    `json:"is_free_tier"`
}

// CheckCredits retrieves the rate limit and credits remaining for the current API key.
func (c *Client) CheckCredits(ctx context.Context) (*KeyData, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.config.BaseURL+"/auth/key", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, &RequestError{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, &RequestError{Err: fmt.Errorf("failed to decode error response: %w", err)}
		}
		return nil, errResp.Error
	}

	var res KeyCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return &res.Data, nil
}
