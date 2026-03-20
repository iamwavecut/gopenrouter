package anthropic

import (
	"encoding/json"
	"net/http"

	"github.com/iamwavecut/gopenrouter/internal/jsonx"
	"github.com/iamwavecut/gopenrouter/internal/sse"
	"github.com/iamwavecut/gopenrouter/shared"
)

type (
	APIError            = shared.APIError
	CacheControl        = shared.CacheControl
	Provider            = shared.Provider
	ProviderPreferences = shared.ProviderPreferences
)

type Request struct {
	Model         string               `json:"model"`
	MaxTokens     int                  `json:"max_tokens,omitempty"`
	Messages      []Message            `json:"messages,omitempty"`
	System        any                  `json:"system,omitempty"`
	Metadata      map[string]any       `json:"metadata,omitempty"`
	StopSequences []string             `json:"stop_sequences,omitempty"`
	Stream        bool                 `json:"stream,omitempty"`
	Temperature   *float64             `json:"temperature,omitempty"`
	TopP          *float64             `json:"top_p,omitempty"`
	TopK          *float64             `json:"top_k,omitempty"`
	Tools         []Tool               `json:"tools,omitempty"`
	ToolChoice    any                  `json:"tool_choice,omitempty"`
	Thinking      any                  `json:"thinking,omitempty"`
	ServiceTier   string               `json:"service_tier,omitempty"`
	Provider      *ProviderPreferences `json:"provider,omitempty"`
	ExtraBody     map[string]any       `json:"-"`
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

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type Tool struct {
	Type           string         `json:"type,omitempty"`
	Name           string         `json:"name,omitempty"`
	Description    string         `json:"description,omitempty"`
	InputSchema    any            `json:"input_schema,omitempty"`
	AllowedDomains []string       `json:"allowed_domains,omitempty"`
	BlockedDomains []string       `json:"blocked_domains,omitempty"`
	MaxUses        *int           `json:"max_uses,omitempty"`
	UserLocation   map[string]any `json:"user_location,omitempty"`
	CacheControl   *CacheControl  `json:"cache_control,omitempty"`
}

type Response struct {
	ID           string         `json:"id,omitempty"`
	Type         string         `json:"type,omitempty"`
	Role         string         `json:"role,omitempty"`
	Model        string         `json:"model,omitempty"`
	Content      []ContentBlock `json:"content,omitempty"`
	StopReason   string         `json:"stop_reason,omitempty"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        map[string]any `json:"usage,omitempty"`
	Raw          map[string]any `json:"-"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	type alias Response
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	*r = Response(decoded)
	r.Raw = raw
	return nil
}

type ContentBlock struct {
	Type string         `json:"type,omitempty"`
	Text string         `json:"text,omitempty"`
	ID   string         `json:"id,omitempty"`
	Name string         `json:"name,omitempty"`
	Raw  map[string]any `json:"-"`
}

func (b *ContentBlock) UnmarshalJSON(data []byte) error {
	type alias ContentBlock
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	*b = ContentBlock(decoded)
	b.Raw = raw
	return nil
}

type Stream struct {
	events *sse.Reader
}

func NewStream(resp *http.Response) *Stream {
	return &Stream{events: sse.NewReader(resp)}
}

type StreamEvent struct {
	Event string
	Type  string         `json:"type,omitempty"`
	Index int            `json:"index,omitempty"`
	Error *APIError      `json:"error,omitempty"`
	Raw   map[string]any `json:"-"`
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
	event.Event = raw.Event
	payload, err := sse.UnwrapPayload(raw.Data)
	if err != nil {
		return event, err
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return event, err
	}
	event.Event = raw.Event
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
