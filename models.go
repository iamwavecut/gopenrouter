package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ModelsList is a response from the /models endpoint.
type ModelsList struct {
	Data []Model `json:"data"`
}

// Model represents a single model from the /models endpoint.
// The schema is inferred from common API patterns.
type Model struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ContextSize  int     `json:"context_length"`
	Architecture any     `json:"architecture"`
	TopProvider  any     `json:"top_provider"`
	Pricing      Pricing `json:"pricing"`
}

// Pricing represents the pricing information for a model.
type Pricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
	Request    string `json:"request"`
	Image      string `json:"image"`
}

// ListModels retrieves the list of available models from OpenRouter.
func (c *Client) ListModels(ctx context.Context) (*ModelsList, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.config.BaseURL+"/models", nil)
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

	var res ModelsList
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return &res, nil
}
