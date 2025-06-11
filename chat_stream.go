package gopenrouter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChatCompletionStreamResponse is a response to a chat completion stream request.
type ChatCompletionStreamResponse struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	Choices           []ChatCompletionStreamChoice `json:"choices"`
	SystemFingerprint string                       `json:"system_fingerprint,omitempty"`
}

// ChatCompletionStreamChoice is a choice in a chat completion stream response.
type ChatCompletionStreamChoice struct {
	Index              int                   `json:"index"`
	Delta              ChatCompletionMessage `json:"delta"`
	FinishReason       string                `json:"finish_reason"`
	NativeFinishReason string                `json:"native_finish_reason,omitempty"`
}

// ChatCompletionStream is a stream of chat completion responses.
type ChatCompletionStream struct {
	*StreamReader
}

// Recv reads a response from the stream.
func (s *ChatCompletionStream) Recv() (response ChatCompletionStreamResponse, err error) {
	buf, err := s.StreamReader.Recv()
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &response)
	return
}

// CreateChatCompletionStream creates a chat completion stream.
func (c *Client) CreateChatCompletionStream(ctx context.Context, r ChatCompletionRequest) (*ChatCompletionStream, error) {
	if err := validateChatCompletionRequest(r); err != nil {
		return nil, err
	}
	r.Stream = true

	req, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+"/chat/completions", r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, &RequestError{Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, &RequestError{Err: fmt.Errorf("failed to decode error response: %w", err)}
		}
		return nil, errResp.Error
	}

	return &ChatCompletionStream{
		StreamReader: &StreamReader{
			reader:   bufio.NewReader(resp.Body),
			response: resp,
		},
	}, nil
}
