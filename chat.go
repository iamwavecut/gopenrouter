package gopenrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
	RoleDeveloper = "developer"
)

type ChatCompletionMessageRole string

type ChatCompletionRequest struct {
	Model         string                  `json:"model,omitempty"`
	Messages      []ChatCompletionMessage `json:"messages,omitempty"`
	Temperature   float64                 `json:"temperature,omitempty"`
	TopP          float64                 `json:"top_p,omitempty"`
	N             int                     `json:"n,omitempty"`
	Stream        bool                    `json:"stream,omitempty"`
	StreamOptions *StreamOptions          `json:"stream_options,omitempty"`
	Stop          []string                `json:"stop,omitempty"`
	// Deprecated: use MaxCompletionTokens instead.
	MaxTokens           int             `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int            `json:"max_completion_tokens,omitempty"`
	PresencePenalty     float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty    float64         `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]int  `json:"logit_bias,omitempty"`
	User                string          `json:"user,omitempty"`
	Seed                *int            `json:"seed,omitempty"`
	Tools               []Tool          `json:"tools,omitempty"`
	ToolChoice          any             `json:"tool_choice,omitempty"`
	LogProbs            *bool           `json:"logprobs,omitempty"`
	ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
	TopK                *float64        `json:"top_k,omitempty"`
	RepetitionPenalty   *float64        `json:"repetition_penalty,omitempty"`
	TopLogProbs         *int            `json:"top_logprobs,omitempty"`
	MinP                *float64        `json:"min_p,omitempty"`
	TopA                *float64        `json:"top_a,omitempty"`
	Prediction          *Prediction     `json:"prediction,omitempty"`

	Models []string `json:"models,omitempty"`
	// Deprecated: use Provider.Sort.Config.Partition instead.
	Route             string               `json:"route,omitempty"`
	Transforms        []string             `json:"transforms,omitempty"`
	Reasoning         *ReasoningParams     `json:"reasoning,omitempty"`
	Usage             *UsageParams         `json:"usage,omitempty"`
	ExtraBody         map[string]any       `json:"-"`
	Plugins           []Plugin             `json:"plugins,omitempty"`
	ParallelToolCalls *bool                `json:"parallel_tool_calls,omitempty"`
	Provider          *ProviderPreferences `json:"provider,omitempty"`
	Metadata          map[string]string    `json:"metadata,omitempty"`
	Modalities        []string             `json:"modalities,omitempty"`
	ImageConfig       map[string]any       `json:"image_config,omitempty"`
	SessionID         string               `json:"session_id,omitempty"`
	Trace             *TraceMetadata       `json:"trace,omitempty"`
	Debug             *DebugOptions        `json:"debug,omitempty"`
}

type chatCompletionRequest ChatCompletionRequest

func (r ChatCompletionRequest) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(chatCompletionRequest(r))
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	maps.Copy(m, r.ExtraBody)
	return json.Marshal(m)
}

type ReasoningEffort string

const (
	ReasoningEffortXHigh   ReasoningEffort = "xhigh"
	ReasoningEffortHigh    ReasoningEffort = "high"
	ReasoningEffortMedium  ReasoningEffort = "medium"
	ReasoningEffortLow     ReasoningEffort = "low"
	ReasoningEffortMinimal ReasoningEffort = "minimal"
	ReasoningEffortNone    ReasoningEffort = "none"
)

type ReasoningSummaryVerbosity string

const (
	ReasoningSummaryAuto     ReasoningSummaryVerbosity = "auto"
	ReasoningSummaryConcise  ReasoningSummaryVerbosity = "concise"
	ReasoningSummaryDetailed ReasoningSummaryVerbosity = "detailed"
)

type ReasoningParams struct {
	Effort  ReasoningEffort           `json:"effort,omitempty"`
	Summary ReasoningSummaryVerbosity `json:"summary,omitempty"`
	// Deprecated: retained only as a legacy compatibility field.
	Exclude bool `json:"exclude,omitempty"`
	// Deprecated: retained only as a legacy compatibility field.
	MaxTokens int `json:"max_tokens,omitempty"`
}

type UsageParams struct {
	Include bool `json:"include"`
}

type ChatCompletionMessage struct {
	Role         ChatCompletionMessageRole   `json:"role"`
	Content      string                      `json:"content,omitempty"`
	MultiContent []ChatCompletionMessagePart `json:"-"`

	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Refusal    string     `json:"refusal,omitempty"`

	Reasoning        string            `json:"reasoning,omitempty"`
	ReasoningDetails []ReasoningDetail `json:"reasoning_details,omitempty"`
	Images           []GeneratedImage  `json:"images,omitempty"`
}

func (m ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"role": m.Role,
	}

	if len(m.MultiContent) > 0 {
		data["content"] = m.MultiContent
	} else if m.Content != "" {
		data["content"] = m.Content
	}
	if m.Name != "" {
		data["name"] = m.Name
	}
	if len(m.ToolCalls) > 0 {
		data["tool_calls"] = m.ToolCalls
	}
	if m.ToolCallID != "" {
		data["tool_call_id"] = m.ToolCallID
	}
	if m.Refusal != "" {
		data["refusal"] = m.Refusal
	}
	if m.Reasoning != "" {
		data["reasoning"] = m.Reasoning
	}
	if len(m.ReasoningDetails) > 0 {
		data["reasoning_details"] = m.ReasoningDetails
	}
	if len(m.Images) > 0 {
		data["images"] = m.Images
	}
	return json.Marshal(data)
}

func (m *ChatCompletionMessage) UnmarshalJSON(data []byte) error {
	type alias ChatCompletionMessage
	var raw struct {
		alias
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*m = ChatCompletionMessage(raw.alias)
	raw.Content = bytes.TrimSpace(raw.Content)
	switch {
	case len(raw.Content) == 0, string(raw.Content) == "null":
		m.Content = ""
		m.MultiContent = nil
	case raw.Content[0] == '[':
		var parts []ChatCompletionMessagePart
		if err := json.Unmarshal(raw.Content, &parts); err != nil {
			return err
		}
		m.MultiContent = parts
		m.Content = ""
	default:
		if err := json.Unmarshal(raw.Content, &m.Content); err != nil {
			return err
		}
		m.MultiContent = nil
	}
	return nil
}

type ChatCompletionMessagePart struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	ImageURL     *ImageURL     `json:"image_url,omitempty"`
	File         *File         `json:"file,omitempty"`
	InputAudio   *InputAudio   `json:"input_audio,omitempty"`
	VideoURL     *VideoInput   `json:"video_url,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
	Arguments   string `json:"arguments,omitempty"`
}

type Prediction struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

type Choice struct {
	Index              int                           `json:"index"`
	Message            ChatCompletionMessage         `json:"message"`
	FinishReason       string                        `json:"finish_reason"`
	NativeFinishReason string                        `json:"native_finish_reason,omitempty"`
	Logprobs           *ChatCompletionChoiceLogprobs `json:"logprobs,omitempty"`
}

type Usage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	Cost                    float64        `json:"cost,omitempty"`
	IsBYOK                  *bool          `json:"is_byok,omitempty"`
	ImageTokens             int            `json:"image_tokens,omitempty"`
	CostDetails             *CostDetails   `json:"cost_details,omitempty"`
	ServerToolUse           *ServerToolUse `json:"server_tool_use,omitempty"`
	PromptTokensDetails     *TokensDetails `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *TokensDetails `json:"completion_tokens_details,omitempty"`
}

type CostDetails struct {
	UpstreamInferenceCost       float64 `json:"upstream_inference_cost,omitempty"`
	UpstreamInferenceInputCost  float64 `json:"upstream_inference_input_cost,omitempty"`
	UpstreamInferenceOutputCost float64 `json:"upstream_inference_output_cost,omitempty"`
}

type ServerToolUse struct {
	WebSearchRequests int `json:"web_search_requests,omitempty"`
}

type TokensDetails struct {
	CachedTokens             int `json:"cached_tokens,omitempty"`
	CacheWriteTokens         int `json:"cache_write_tokens,omitempty"`
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	AudioTokens              int `json:"audio_tokens,omitempty"`
	VideoTokens              int `json:"video_tokens,omitempty"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

type ChatCompletionChoiceLogprobs struct {
	Content []ChatCompletionTokenLogprob `json:"content,omitempty"`
	Refusal []ChatCompletionTokenLogprob `json:"refusal,omitempty"`
}

type ChatCompletionTokenLogprob struct {
	Token       string                                 `json:"token"`
	LogProb     float64                                `json:"logprob"`
	Bytes       []byte                                 `json:"bytes,omitempty"`
	TopLogProbs []ChatCompletionTokenLogprobTopLogprob `json:"top_logprobs,omitempty"`
}

type ChatCompletionTokenLogprobTopLogprob struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

type ReasoningDetail struct {
	Type      string `json:"type"`
	Summary   string `json:"summary,omitempty"`
	Data      string `json:"data,omitempty"`
	Text      string `json:"text,omitempty"`
	Signature string `json:"signature,omitempty"`
	ID        string `json:"id,omitempty"`
	Format    string `json:"format,omitempty"`
	Index     *int   `json:"index,omitempty"`

	Encrypted string `json:"-"`
}

func (r ReasoningDetail) MarshalJSON() ([]byte, error) {
	type alias ReasoningDetail
	if r.Data == "" && r.Encrypted != "" {
		r.Data = r.Encrypted
	}
	return json.Marshal(alias(r))
}

func (r *ReasoningDetail) UnmarshalJSON(data []byte) error {
	type alias ReasoningDetail
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*r = ReasoningDetail(decoded)
	if r.Data != "" {
		r.Encrypted = r.Data
	}
	return nil
}

type ToolCall struct {
	Index    int      `json:"index,omitempty"`
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

func (c *Client) CreateChatCompletion(ctx context.Context, r ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if err := validateChatCompletionRequest(r); err != nil {
		return nil, err
	}
	if r.Stream {
		return nil, &RequestError{Err: fmt.Errorf("use CreateChatCompletionStream for streaming requests")}
	}

	var res ChatCompletionResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/chat/completions", nil, r, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func validateChatCompletionRequest(r ChatCompletionRequest) error {
	if r.MaxTokens > 0 && r.MaxCompletionTokens != nil && *r.MaxCompletionTokens > 0 {
		return &RequestError{Err: fmt.Errorf("only one of 'max_tokens' and 'max_completion_tokens' can be specified")}
	}
	if r.Route != "" && r.Provider != nil && r.Provider.Sort != nil && r.Provider.Sort.Config != nil && r.Provider.Sort.Config.Partition != "" {
		return &RequestError{Err: fmt.Errorf("only one of 'route' and 'provider.sort.partition' can be specified")}
	}
	return nil
}
