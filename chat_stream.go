package gopenrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

type ChatCompletionStreamResponse struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	Choices           []ChatCompletionStreamChoice `json:"choices"`
	SystemFingerprint string                       `json:"system_fingerprint,omitempty"`
	Usage             *Usage                       `json:"usage,omitempty"`
	Error             *APIError                    `json:"error,omitempty"`
}

type ChatCompletionStreamChoice struct {
	Index              int                                 `json:"index"`
	Delta              ChatCompletionMessage               `json:"delta"`
	FinishReason       string                              `json:"finish_reason"`
	NativeFinishReason string                              `json:"native_finish_reason,omitempty"`
	Logprobs           *ChatCompletionStreamChoiceLogprobs `json:"logprobs,omitempty"`
}

type ChatCompletionStreamChoiceLogprobs struct {
	Content []ChatCompletionStreamChoiceDelta `json:"content,omitempty"`
	Refusal []ChatCompletionStreamChoiceDelta `json:"refusal,omitempty"`
}

type ChatCompletionStreamChoiceDelta struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

type ChatCompletionStream struct {
	*StreamReader
}

func (s *ChatCompletionStream) Recv() (ChatCompletionStreamResponse, error) {
	var response ChatCompletionStreamResponse
	buf, err := s.StreamReader.Recv()
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

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
		defer resp.Body.Close()
		return nil, decodeErrorResponse(resp)
	}

	return &ChatCompletionStream{
		StreamReader: newStreamReader(resp),
	}, nil
}
