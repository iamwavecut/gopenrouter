package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"strconv"

	"github.com/iamwavecut/gopenrouter/internal/jsonx"
)

type BigNumber string

func (n BigNumber) MarshalJSON() ([]byte, error) {
	raw := string(n)
	if raw == "" {
		return []byte("null"), nil
	}
	if _, err := strconv.ParseFloat(raw, 64); err == nil {
		return []byte(raw), nil
	}
	return json.Marshal(raw)
}

func (n *BigNumber) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*n = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*n = BigNumber(s)
		return nil
	}

	*n = BigNumber(string(data))
	return nil
}

type DataCollection string

const (
	DataCollectionAllow DataCollection = "allow"
	DataCollectionDeny  DataCollection = "deny"
)

type Quantization string

type ProviderSort string

const (
	ProviderSortPrice      ProviderSort = "price"
	ProviderSortThroughput ProviderSort = "throughput"
	ProviderSortLatency    ProviderSort = "latency"
)

type ProviderSortPartition string

const (
	ProviderSortPartitionModel ProviderSortPartition = "model"
	ProviderSortPartitionNone  ProviderSortPartition = "none"
)

type ProviderSortConfig struct {
	By        ProviderSort          `json:"by,omitempty"`
	Partition ProviderSortPartition `json:"partition,omitempty"`
}

type ProviderSortPreference struct {
	Mode   ProviderSort
	Config *ProviderSortConfig
}

func (p ProviderSortPreference) MarshalJSON() ([]byte, error) {
	if p.Config != nil {
		return json.Marshal(p.Config)
	}
	if p.Mode != "" {
		return json.Marshal(p.Mode)
	}
	return []byte("null"), nil
}

func (p *ProviderSortPreference) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*p = ProviderSortPreference{}
		return nil
	}

	var mode ProviderSort
	if err := json.Unmarshal(data, &mode); err == nil {
		p.Mode = mode
		p.Config = nil
		return nil
	}

	var cfg ProviderSortConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	p.Mode = ""
	p.Config = &cfg
	return nil
}

type PercentileStats struct {
	P50 *float64 `json:"p50,omitempty"`
	P75 *float64 `json:"p75,omitempty"`
	P90 *float64 `json:"p90,omitempty"`
	P99 *float64 `json:"p99,omitempty"`
}

type (
	ThroughputCutoffs = PercentileStats
	LatencyCutoffs    = PercentileStats
)

type ThroughputPreference struct {
	Value   *float64
	Cutoffs *ThroughputCutoffs
}

func (p ThroughputPreference) MarshalJSON() ([]byte, error) {
	if p.Cutoffs != nil {
		return json.Marshal(p.Cutoffs)
	}
	if p.Value != nil {
		return json.Marshal(*p.Value)
	}
	return []byte("null"), nil
}

func (p *ThroughputPreference) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*p = ThroughputPreference{}
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err == nil {
		p.Value = &value
		p.Cutoffs = nil
		return nil
	}

	var cutoffs ThroughputCutoffs
	if err := json.Unmarshal(data, &cutoffs); err != nil {
		return err
	}
	p.Value = nil
	p.Cutoffs = &cutoffs
	return nil
}

type LatencyPreference struct {
	Value   *float64
	Cutoffs *LatencyCutoffs
}

func (p LatencyPreference) MarshalJSON() ([]byte, error) {
	if p.Cutoffs != nil {
		return json.Marshal(p.Cutoffs)
	}
	if p.Value != nil {
		return json.Marshal(*p.Value)
	}
	return []byte("null"), nil
}

func (p *LatencyPreference) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*p = LatencyPreference{}
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err == nil {
		p.Value = &value
		p.Cutoffs = nil
		return nil
	}

	var cutoffs LatencyCutoffs
	if err := json.Unmarshal(data, &cutoffs); err != nil {
		return err
	}
	p.Value = nil
	p.Cutoffs = &cutoffs
	return nil
}

type ProviderMaxPrice struct {
	Prompt     BigNumber `json:"prompt,omitempty"`
	Completion BigNumber `json:"completion,omitempty"`
	Image      BigNumber `json:"image,omitempty"`
	Audio      BigNumber `json:"audio,omitempty"`
	Request    BigNumber `json:"request,omitempty"`
}

type ProviderPreferences struct {
	AllowFallbacks         *bool                   `json:"allow_fallbacks,omitempty"`
	RequireParameters      *bool                   `json:"require_parameters,omitempty"`
	DataCollection         DataCollection          `json:"data_collection,omitempty"`
	ZDR                    *bool                   `json:"zdr,omitempty"`
	EnforceDistillableText *bool                   `json:"enforce_distillable_text,omitempty"`
	Order                  []string                `json:"order,omitempty"`
	Only                   []string                `json:"only,omitempty"`
	Ignore                 []string                `json:"ignore,omitempty"`
	Quantizations          []Quantization          `json:"quantizations,omitempty"`
	Sort                   *ProviderSortPreference `json:"sort,omitempty"`
	MaxPrice               *ProviderMaxPrice       `json:"max_price,omitempty"`
	PreferredMinThroughput *ThroughputPreference   `json:"preferred_min_throughput,omitempty"`
	PreferredMaxLatency    *LatencyPreference      `json:"preferred_max_latency,omitempty"`
}

type Provider = ProviderPreferences

type TraceMetadata struct {
	TraceID        string         `json:"trace_id,omitempty"`
	TraceName      string         `json:"trace_name,omitempty"`
	SpanName       string         `json:"span_name,omitempty"`
	GenerationName string         `json:"generation_name,omitempty"`
	ParentSpanID   string         `json:"parent_span_id,omitempty"`
	Extra          map[string]any `json:"-"`
}

func (t TraceMetadata) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	if t.TraceID != "" {
		m["trace_id"] = t.TraceID
	}
	if t.TraceName != "" {
		m["trace_name"] = t.TraceName
	}
	if t.SpanName != "" {
		m["span_name"] = t.SpanName
	}
	if t.GenerationName != "" {
		m["generation_name"] = t.GenerationName
	}
	if t.ParentSpanID != "" {
		m["parent_span_id"] = t.ParentSpanID
	}
	maps.Copy(m, t.Extra)
	return json.Marshal(m)
}

func (t *TraceMetadata) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	type alias TraceMetadata
	var known alias
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	delete(raw, "trace_id")
	delete(raw, "trace_name")
	delete(raw, "span_name")
	delete(raw, "generation_name")
	delete(raw, "parent_span_id")
	*t = TraceMetadata(known)
	if len(raw) > 0 {
		t.Extra = raw
	}
	return nil
}

type DebugOptions struct {
	EchoUpstreamBody bool `json:"echo_upstream_body,omitempty"`
}

type CacheControl struct {
	Type string `json:"type"`
	TTL  string `json:"ttl,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type File struct {
	Filename string `json:"filename,omitempty"`
	FileData string `json:"file_data,omitempty"`
}

type InputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

type VideoInput struct {
	URL string `json:"url"`
}

type GeneratedImage struct {
	ImageURL ImageURL `json:"image_url"`
}

type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
	Grammar    string      `json:"grammar,omitempty"`
}

type JSONSchema struct {
	Name   string `json:"name"`
	Strict bool   `json:"strict,omitempty"`
	Schema any    `json:"schema"`
}

type Plugin struct {
	ID            PluginID       `json:"id"`
	Enabled       *bool          `json:"enabled,omitempty"`
	AllowedModels []string       `json:"allowed_models,omitempty"`
	MaxResults    *float64       `json:"max_results,omitempty"`
	SearchPrompt  string         `json:"search_prompt,omitempty"`
	Engine        string         `json:"engine,omitempty"`
	PDF           *PDFPlugin     `json:"pdf,omitempty"`
	Config        any            `json:"-"`
	Extra         map[string]any `json:"-"`
}

func (p Plugin) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"id": p.ID,
	}
	if p.Enabled != nil {
		data["enabled"] = *p.Enabled
	}
	if len(p.AllowedModels) > 0 {
		data["allowed_models"] = p.AllowedModels
	}
	if p.MaxResults != nil {
		data["max_results"] = *p.MaxResults
	}
	if p.SearchPrompt != "" {
		data["search_prompt"] = p.SearchPrompt
	}
	if p.Engine != "" {
		data["engine"] = p.Engine
	}
	if p.PDF != nil {
		data["pdf"] = p.PDF
	}
	if p.Config != nil {
		m, err := jsonx.MarshalMap(p.Config)
		if err != nil {
			return nil, err
		}
		jsonx.MergeMaps(data, m)
	}
	jsonx.MergeMaps(data, p.Extra)
	return json.Marshal(data)
}

func (p *Plugin) UnmarshalJSON(data []byte) error {
	type alias Plugin
	var known alias
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	raw, err := jsonx.MarshalMap(json.RawMessage(data))
	if err != nil {
		return err
	}
	delete(raw, "id")
	delete(raw, "enabled")
	delete(raw, "allowed_models")
	delete(raw, "max_results")
	delete(raw, "search_prompt")
	delete(raw, "engine")
	delete(raw, "pdf")

	*p = Plugin(known)
	if len(raw) > 0 {
		p.Extra = raw
	}
	return nil
}

type PluginID string

const (
	PluginIDAutoRouter      PluginID = "auto-router"
	PluginIDModeration      PluginID = "moderation"
	PluginIDFileParser      PluginID = "file-parser"
	PluginIDWeb             PluginID = "web"
	PluginIDResponseHealing PluginID = "response-healing"
)

type PDFEngine string

const (
	PDFEngineMistralOCR PDFEngine = "mistral-ocr"
	PDFEnginePDFText    PDFEngine = "pdf-text"
	PDFEngineNative     PDFEngine = "native"
)

type PDFPlugin struct {
	Engine string `json:"engine,omitempty"`
}

type FileParserConfig struct {
	PDF *PDFPlugin `json:"pdf,omitempty"`
}

type WebSearchOptions struct {
	SearchContextSize SearchContextSize `json:"search_context_size,omitempty"`
	MaxResults        *float64          `json:"max_results,omitempty"`
	SearchPrompt      string            `json:"search_prompt,omitempty"`
	Engine            string            `json:"engine,omitempty"`
}

type SearchContextSize string

const (
	SearchContextSizeLow    SearchContextSize = "low"
	SearchContextSizeMedium SearchContextSize = "medium"
	SearchContextSizeHigh   SearchContextSize = "high"
)

type ErrorResponse struct {
	Error  *APIError `json:"error,omitempty"`
	Type   string    `json:"type,omitempty"`
	UserID string    `json:"user_id,omitempty"`
}

type ProviderError map[string]any

func (e *ProviderError) Message() any {
	if e == nil {
		return nil
	}
	if msg, ok := (*e)["message"]; ok {
		return msg
	}
	return nil
}

type APIError struct {
	Type          string         `json:"type,omitempty"`
	Message       string         `json:"message"`
	Code          any            `json:"code"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	ProviderError *ProviderError `json:"provider_error,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

type RequestError struct {
	HTTPStatus     string
	HTTPStatusCode int
	Err            error
	Body           []byte
}

func (e *RequestError) Error() string {
	if e == nil {
		return "request error: <nil>"
	}
	if e.HTTPStatus != "" {
		return fmt.Sprintf("request error: %s: %v", e.HTTPStatus, e.Err)
	}
	return fmt.Sprintf("request error: %v", e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
