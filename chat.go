package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
)

type ChatCompletionMessageRole string

// ChatCompletionRequest is a request to create a chat completion.
type ChatCompletionRequest struct {
	Model             string                  `json:"model"`
	Messages          []ChatCompletionMessage `json:"messages"`
	Temperature       float64                 `json:"temperature,omitempty"`
	TopP              float64                 `json:"top_p,omitempty"`
	N                 int                     `json:"n,omitempty"`
	Stream            bool                    `json:"stream,omitempty"`
	Stop              []string                `json:"stop,omitempty"`
	MaxTokens         int                     `json:"max_tokens,omitempty"`
	PresencePenalty   float64                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty  float64                 `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int          `json:"logit_bias,omitempty"`
	User              string                  `json:"user,omitempty"`
	Seed              *int                    `json:"seed,omitempty"`
	Tools             []Tool                  `json:"tools,omitempty"`
	ToolChoice        any                     `json:"tool_choice,omitempty"`
	ResponseFormat    *ResponseFormat         `json:"response_format,omitempty"`
	TopK              *float64                `json:"top_k,omitempty"`
	RepetitionPenalty *float64                `json:"repetition_penalty,omitempty"`
	TopLogProbs       *int                    `json:"top_logprobs,omitempty"`
	MinP              *float64                `json:"min_p,omitempty"`
	TopA              *float64                `json:"top_a,omitempty"`
	Prediction        *Prediction             `json:"prediction,omitempty"`

	// OpenRouter-specific fields
	Models     []string         `json:"models,omitempty"`
	Route      string           `json:"route,omitempty"`
	Transforms []string         `json:"transforms,omitempty"`
	Reasoning  *ReasoningParams `json:"reasoning,omitempty"`
	Usage      *UsageParams     `json:"usage,omitempty"`
	ExtraBody  map[string]any   `json:"-"`
}

// chatCompletionRequest is an alias for ChatCompletionRequest to avoid recursion in MarshalJSON.
type chatCompletionRequest ChatCompletionRequest

// MarshalJSON implements json.Marshaler.
func (r ChatCompletionRequest) MarshalJSON() ([]byte, error) {
	// Use an alias to avoid recursive call to MarshalJSON
	b, err := json.Marshal(chatCompletionRequest(r))
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	// Merge ExtraBody
	for k, v := range r.ExtraBody {
		m[k] = v
	}

	return json.Marshal(m)
}

// ReasoningParams specifies parameters for reasoning token generation.
// One of MaxTokens or Effort should be specified.
type ReasoningParams struct {
	Exclude   bool   `json:"exclude,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty"`
	Effort    string `json:"effort,omitempty"` // Can be "high", "medium", "low"
}

// UsageParams is for OpenRouter's usage accounting feature.
type UsageParams struct {
	Include bool `json:"include"`
}

// ChatCompletionMessage is a message in a chat completion request.
type ChatCompletionMessage struct {
	Role         ChatCompletionMessageRole   `json:"role"`
	Content      string                      `json:"content,omitempty"`
	MultiContent []ChatCompletionMessagePart `json:"-"` // Will be marshalled to content

	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`

	// For response messages
	Reasoning        string            `json:"reasoning,omitempty"`
	ReasoningDetails []ReasoningDetail `json:"reasoning_details,omitempty"`
}

// MarshalJSON handles marshalling the Content or MultiContent field.
func (m *ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	type alias struct {
		Role       ChatCompletionMessageRole `json:"role"`
		Name       string                    `json:"name,omitempty"`
		ToolCalls  []ToolCall                `json:"tool_calls,omitempty"`
		ToolCallID string                    `json:"tool_call_id,omitempty"`
	}

	if len(m.MultiContent) > 0 {
		// If MultiContent is present, marshal it as 'content' and omit the original string content
		return json.Marshal(&struct {
			alias
			Content []ChatCompletionMessagePart `json:"content"`
		}{
			alias: alias{
				Role:       m.Role,
				Name:       m.Name,
				ToolCalls:  m.ToolCalls,
				ToolCallID: m.ToolCallID,
			},
			Content: m.MultiContent,
		})
	}

	// Otherwise, use the default marshalling which includes the string 'Content'
	type defaultAlias ChatCompletionMessage
	return json.Marshal((*defaultAlias)(m))
}

// ChatCompletionMessagePart is a part of a multi-part message.
type ChatCompletionMessagePart struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	ImageURL     *ImageURL     `json:"image_url,omitempty"`
	File         *File         `json:"file,omitempty"` // OpenRouter specific
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// CacheControl is for controlling caching of message parts.
type CacheControl struct {
	Type string `json:"type"` // e.g. "ephemeral"
}

// ImageURL is the URL of an image.
type ImageURL struct {
	URL string `json:"url"`
}

// File is a file attachment. OpenRouter specific.
type File struct {
	Filename string `json:"filename"`
	FileData string `json:"file_data"` // base64 encoded
}

// ResponseFormat is the format of the response.
type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

// JSONSchema is the schema for the JSON response format.
type JSONSchema struct {
	Name   string `json:"name"`
	Strict bool   `json:"strict"`
	Schema any    `json:"schema"`
}

// Provider is the provider information.
type Provider struct{}

// Tool is a tool the model can use.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function is a function the model can call.
type Function struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
	Arguments   string `json:"arguments,omitempty"` // For tool call responses
}

// ChatCompletionResponse is a response to a chat completion request.
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// Choice is a choice in a chat completion response.
type Choice struct {
	Index              int                   `json:"index"`
	Message            ChatCompletionMessage `json:"message"`
	FinishReason       string                `json:"finish_reason"`
	NativeFinishReason string                `json:"native_finish_reason,omitempty"`
}

// Usage is the usage of the model.
type Usage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	Cost                    float64        `json:"cost,omitempty"`
	PromptTokensDetails     *TokensDetails `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *TokensDetails `json:"completion_tokens_details,omitempty"`
}

// TokensDetails provides more detail on token usage.
type TokensDetails struct {
	CachedTokens    int `json:"cached_tokens,omitempty"`
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// ReasoningDetail is a detail about the reasoning process.
type ReasoningDetail struct {
	Type      string `json:"type"`
	Summary   string `json:"summary,omitempty"`
	Encrypted string `json:"encrypted,omitempty"`
	Text      string `json:"text,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// ToolCall is a call to a tool.
type ToolCall struct {
	Index    int      `json:"index,omitempty"`
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

// Prediction is the predicted output for latency optimization.
type Prediction struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// CreateChatCompletion creates a chat completion.
func (c *Client) CreateChatCompletion(ctx context.Context, r ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if err := validateChatCompletionRequest(r); err != nil {
		return nil, err
	}

	if r.Stream {
		return nil, &RequestError{Err: fmt.Errorf("use CreateChatCompletionStream for streaming requests")}
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+"/chat/completions", r)
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

	var res ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, &RequestError{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return &res, nil
}

func validateChatCompletionRequest(r ChatCompletionRequest) error {
	if r.Reasoning != nil && r.Reasoning.MaxTokens > 0 && r.Reasoning.Effort != "" {
		return &RequestError{Err: fmt.Errorf("only one of 'reasoning.effort' and 'reasoning.max_tokens' can be specified")}
	}
	return nil
}
