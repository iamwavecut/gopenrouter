package catalog

import "github.com/iamwavecut/gopenrouter/shared"

type (
	BigNumber       = shared.BigNumber
	Quantization    = shared.Quantization
	PercentileStats = shared.PercentileStats
)

type ModelsList struct {
	Data []Model `json:"data"`
}

type Model struct {
	ID                  string             `json:"id"`
	CanonicalSlug       string             `json:"canonical_slug,omitempty"`
	HuggingFaceID       string             `json:"hugging_face_id,omitempty"`
	Name                string             `json:"name"`
	Created             int64              `json:"created,omitempty"`
	Description         string             `json:"description,omitempty"`
	Pricing             Pricing            `json:"pricing"`
	ContextSize         int                `json:"context_length,omitempty"`
	Architecture        *ModelArchitecture `json:"architecture,omitempty"`
	TopProvider         *TopProviderInfo   `json:"top_provider,omitempty"`
	PerRequestLimits    *PerRequestLimits  `json:"per_request_limits,omitempty"`
	SupportedParameters []string           `json:"supported_parameters,omitempty"`
	DefaultParameters   *DefaultParameters `json:"default_parameters,omitempty"`
	ExpirationDate      string             `json:"expiration_date,omitempty"`
}

type Pricing struct {
	Prompt            BigNumber `json:"prompt,omitempty"`
	Completion        BigNumber `json:"completion,omitempty"`
	Request           BigNumber `json:"request,omitempty"`
	Image             BigNumber `json:"image,omitempty"`
	ImageToken        BigNumber `json:"image_token,omitempty"`
	ImageOutput       BigNumber `json:"image_output,omitempty"`
	Audio             BigNumber `json:"audio,omitempty"`
	AudioOutput       BigNumber `json:"audio_output,omitempty"`
	InputAudioCache   BigNumber `json:"input_audio_cache,omitempty"`
	WebSearch         BigNumber `json:"web_search,omitempty"`
	InternalReasoning BigNumber `json:"internal_reasoning,omitempty"`
	InputCacheRead    BigNumber `json:"input_cache_read,omitempty"`
	InputCacheWrite   BigNumber `json:"input_cache_write,omitempty"`
	Discount          float64   `json:"discount,omitempty"`
}

type ModelArchitecture struct {
	Tokenizer        string   `json:"tokenizer,omitempty"`
	InstructType     string   `json:"instruct_type,omitempty"`
	Modality         string   `json:"modality,omitempty"`
	InputModalities  []string `json:"input_modalities,omitempty"`
	OutputModalities []string `json:"output_modalities,omitempty"`
}

type TopProviderInfo struct {
	ContextLength       int  `json:"context_length,omitempty"`
	MaxCompletionTokens int  `json:"max_completion_tokens,omitempty"`
	IsModerated         bool `json:"is_moderated"`
}

type PerRequestLimits struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
}

type DefaultParameters struct {
	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"top_p,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
}

type ModelsListParams struct {
	Category            string `json:"category,omitempty"`
	SupportedParameters string `json:"supported_parameters,omitempty"`
}

type ModelsCountResponse struct {
	Data struct {
		Count int `json:"count"`
	} `json:"data"`
}

type ProvidersList struct {
	Data []ProviderInfo `json:"data"`
}

type ProviderInfo struct {
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	PrivacyPolicyURL  string `json:"privacy_policy_url,omitempty"`
	TermsOfServiceURL string `json:"terms_of_service_url,omitempty"`
	StatusPageURL     string `json:"status_page_url,omitempty"`
}

type PublicEndpoint struct {
	Name                    string           `json:"name"`
	ModelID                 string           `json:"model_id"`
	ModelName               string           `json:"model_name"`
	ContextLength           int              `json:"context_length,omitempty"`
	Pricing                 Pricing          `json:"pricing"`
	ProviderName            string           `json:"provider_name,omitempty"`
	Tag                     string           `json:"tag,omitempty"`
	Quantization            Quantization     `json:"quantization,omitempty"`
	MaxCompletionTokens     int              `json:"max_completion_tokens,omitempty"`
	MaxPromptTokens         int              `json:"max_prompt_tokens,omitempty"`
	SupportedParameters     []string         `json:"supported_parameters,omitempty"`
	Status                  any              `json:"status,omitempty"`
	UptimeLast30M           *float64         `json:"uptime_last_30m,omitempty"`
	SupportsImplicitCaching bool             `json:"supports_implicit_caching,omitempty"`
	LatencyLast30M          *PercentileStats `json:"latency_last_30m,omitempty"`
	ThroughputLast30M       *PercentileStats `json:"throughput_last_30m,omitempty"`
}

type ZDREndpointsList struct {
	Data []PublicEndpoint `json:"data"`
}

type ModelEndpoints struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Created      int64              `json:"created,omitempty"`
	Description  string             `json:"description,omitempty"`
	Architecture *ModelArchitecture `json:"architecture,omitempty"`
	Endpoints    []PublicEndpoint   `json:"endpoints"`
}

type ModelEndpointsResponse struct {
	Data ModelEndpoints `json:"data"`
}
