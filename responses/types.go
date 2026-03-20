package responses

import (
	"encoding/json"
	"net/http"

	"github.com/iamwavecut/gopenrouter/internal/jsonx"
	"github.com/iamwavecut/gopenrouter/internal/sse"
	"github.com/iamwavecut/gopenrouter/shared"
)

type (
	Plugin              = shared.Plugin
	Provider            = shared.Provider
	ProviderPreferences = shared.ProviderPreferences
	ResponseFormat      = shared.ResponseFormat
	TraceMetadata       = shared.TraceMetadata
	APIError            = shared.APIError
)

type Request struct {
	Input              any                  `json:"input,omitempty"`
	Instructions       string               `json:"instructions,omitempty"`
	Metadata           map[string]string    `json:"metadata,omitempty"`
	Tools              []Tool               `json:"tools,omitempty"`
	ToolChoice         any                  `json:"tool_choice,omitempty"`
	ParallelToolCalls  *bool                `json:"parallel_tool_calls,omitempty"`
	Model              string               `json:"model,omitempty"`
	Models             []string             `json:"models,omitempty"`
	Text               *TextConfig          `json:"text,omitempty"`
	Reasoning          *ReasoningConfig     `json:"reasoning,omitempty"`
	MaxOutputTokens    *int                 `json:"max_output_tokens,omitempty"`
	Temperature        *float64             `json:"temperature,omitempty"`
	TopP               *float64             `json:"top_p,omitempty"`
	TopLogProbs        *int                 `json:"top_logprobs,omitempty"`
	MaxToolCalls       *int                 `json:"max_tool_calls,omitempty"`
	PresencePenalty    *float64             `json:"presence_penalty,omitempty"`
	FrequencyPenalty   *float64             `json:"frequency_penalty,omitempty"`
	TopK               *float64             `json:"top_k,omitempty"`
	ImageConfig        map[string]any       `json:"image_config,omitempty"`
	Modalities         []string             `json:"modalities,omitempty"`
	PromptCacheKey     string               `json:"prompt_cache_key,omitempty"`
	PreviousResponseID string               `json:"previous_response_id,omitempty"`
	Prompt             any                  `json:"prompt,omitempty"`
	Include            []string             `json:"include,omitempty"`
	Background         *bool                `json:"background,omitempty"`
	SafetyIdentifier   string               `json:"safety_identifier,omitempty"`
	Store              *bool                `json:"store,omitempty"`
	ServiceTier        string               `json:"service_tier,omitempty"`
	Truncation         any                  `json:"truncation,omitempty"`
	Stream             bool                 `json:"stream,omitempty"`
	Provider           *ProviderPreferences `json:"provider,omitempty"`
	Plugins            []Plugin             `json:"plugins,omitempty"`
	Route              string               `json:"route,omitempty"`
	User               string               `json:"user,omitempty"`
	SessionID          string               `json:"session_id,omitempty"`
	Trace              *TraceMetadata       `json:"trace,omitempty"`
	ExtraBody          map[string]any       `json:"-"`
}

type requestAlias Request

func (r Request) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(requestAlias(r))
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	jsonx.MergeMaps(m, r.ExtraBody)
	return json.Marshal(m)
}

type Tool struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

type TextConfig struct {
	Format *ResponseFormat `json:"format,omitempty"`
}

type ReasoningConfig struct {
	Effort    string `json:"effort,omitempty"`
	Summary   string `json:"summary,omitempty"`
	MaxTokens *int   `json:"max_tokens,omitempty"`
	Enabled   *bool  `json:"enabled,omitempty"`
}

type Response struct {
	ID                string            `json:"id"`
	Object            string            `json:"object,omitempty"`
	CreatedAt         int64             `json:"created_at,omitempty"`
	Model             string            `json:"model,omitempty"`
	Status            string            `json:"status,omitempty"`
	Output            []OutputItem      `json:"output,omitempty"`
	Usage             *Usage            `json:"usage,omitempty"`
	Tools             []Tool            `json:"tools,omitempty"`
	ToolChoice        any               `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool             `json:"parallel_tool_calls,omitempty"`
	Error             *APIError         `json:"error,omitempty"`
	IncompleteDetails map[string]any    `json:"incomplete_details,omitempty"`
	Temperature       *float64          `json:"temperature,omitempty"`
	TopP              *float64          `json:"top_p,omitempty"`
	MaxOutputTokens   *int              `json:"max_output_tokens,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Instructions      string            `json:"instructions,omitempty"`
}

type Usage struct {
	InputTokens         int            `json:"input_tokens"`
	OutputTokens        int            `json:"output_tokens"`
	TotalTokens         int            `json:"total_tokens"`
	InputTokensDetails  map[string]any `json:"input_tokens_details,omitempty"`
	OutputTokensDetails map[string]any `json:"output_tokens_details,omitempty"`
	Cost                float64        `json:"cost,omitempty"`
	IsBYOK              *bool          `json:"is_byok,omitempty"`
	CostDetails         map[string]any `json:"cost_details,omitempty"`
}

type OutputItem struct {
	ID        string         `json:"id,omitempty"`
	Type      string         `json:"type,omitempty"`
	Status    string         `json:"status,omitempty"`
	Role      string         `json:"role,omitempty"`
	Name      string         `json:"name,omitempty"`
	Arguments string         `json:"arguments,omitempty"`
	CallID    string         `json:"call_id,omitempty"`
	Content   []ContentPart  `json:"content,omitempty"`
	Summary   []ContentPart  `json:"summary,omitempty"`
	Raw       map[string]any `json:"-"`
}

func (i *OutputItem) UnmarshalJSON(data []byte) error {
	type alias OutputItem
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	*i = OutputItem(decoded)
	i.Raw = raw
	return nil
}

type ContentPart struct {
	Type        string           `json:"type,omitempty"`
	Text        string           `json:"text,omitempty"`
	Annotations []map[string]any `json:"annotations,omitempty"`
	Raw         map[string]any   `json:"-"`
}

func (p *ContentPart) UnmarshalJSON(data []byte) error {
	type alias ContentPart
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	*p = ContentPart(decoded)
	p.Raw = raw
	return nil
}

type Stream struct {
	events *sse.Reader
}

func NewStream(resp *http.Response) *Stream {
	return &Stream{events: sse.NewReader(resp)}
}

type StreamEvent struct {
	Type           string         `json:"type"`
	SequenceNumber int            `json:"sequence_number,omitempty"`
	OutputIndex    int            `json:"output_index,omitempty"`
	ItemID         string         `json:"item_id,omitempty"`
	ContentIndex   int            `json:"content_index,omitempty"`
	Delta          string         `json:"delta,omitempty"`
	Text           string         `json:"text,omitempty"`
	Name           string         `json:"name,omitempty"`
	Arguments      string         `json:"arguments,omitempty"`
	Response       *Response      `json:"response,omitempty"`
	Item           *OutputItem    `json:"item,omitempty"`
	Part           *ContentPart   `json:"part,omitempty"`
	Error          *APIError      `json:"error,omitempty"`
	Raw            map[string]any `json:"-"`
}

func (e *StreamEvent) UnmarshalJSON(data []byte) error {
	type alias StreamEvent
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	*e = StreamEvent(decoded)
	e.Raw = raw
	return nil
}

func (s *Stream) Recv() (StreamEvent, error) {
	var event StreamEvent
	raw, err := s.events.RecvEvent()
	if err != nil {
		return event, err
	}
	payload, err := sse.UnwrapPayload(raw.Data)
	if err != nil {
		return event, err
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return event, err
	}
	if event.Error != nil {
		return event, event.Error
	}
	return event, nil
}

func (s *Stream) Close() {
	if s == nil || s.events == nil {
		return
	}
	s.events.Close()
}
