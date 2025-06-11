package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GenerationResponse is the top-level response from the /generation endpoint.
type GenerationResponse struct {
	Data Generation `json:"data"`
}

// Generation contains the detailed statistics for a single generation.
type Generation struct {
	ID                     string  `json:"id"`
	Model                  string  `json:"model"`
	CreatedAt              string  `json:"created_at"`
	PromptTokens           int     `json:"tokens_prompt"`
	CompletionTokens       int     `json:"tokens_completion"`
	NativePromptTokens     int     `json:"native_tokens_prompt"`
	NativeCompletionTokens int     `json:"native_tokens_completion"`
	FinishReason           string  `json:"finish_reason"`
	NativeFinishReason     string  `json:"native_finish_reason"`
	TotalCost              float64 `json:"total_cost"`
}

// GetGeneration retrieves the full generation information for a given ID.
// This is useful for getting precise token counts and cost for a request.
func (c *Client) GetGeneration(ctx context.Context, id string) (*Generation, error) {
	req, err := c.newRequest(ctx, http.MethodGet, c.config.BaseURL+"/generation?id="+id, nil)
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

	var res GenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return &res.Data, nil
}
