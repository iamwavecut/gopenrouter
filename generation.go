package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
)

type GenerationResponse struct {
	Data Generation `json:"data"`
}

type Generation struct {
	ID                          string                       `json:"id"`
	UpstreamID                  string                       `json:"upstream_id,omitempty"`
	TotalCost                   float64                      `json:"total_cost,omitempty"`
	CacheDiscount               float64                      `json:"cache_discount,omitempty"`
	UpstreamInferenceCost       float64                      `json:"upstream_inference_cost,omitempty"`
	CreatedAt                   string                       `json:"created_at,omitempty"`
	Model                       string                       `json:"model,omitempty"`
	AppID                       int                          `json:"app_id,omitempty"`
	Streamed                    *bool                        `json:"streamed,omitempty"`
	Cancelled                   *bool                        `json:"cancelled,omitempty"`
	ProviderName                string                       `json:"provider_name,omitempty"`
	Latency                     *float64                     `json:"latency,omitempty"`
	ModerationLatency           *float64                     `json:"moderation_latency,omitempty"`
	GenerationTime              *float64                     `json:"generation_time,omitempty"`
	FinishReason                string                       `json:"finish_reason,omitempty"`
	NativeFinishReason          string                       `json:"native_finish_reason,omitempty"`
	PromptTokens                int                          `json:"tokens_prompt,omitempty"`
	CompletionTokens            int                          `json:"tokens_completion,omitempty"`
	NativePromptTokens          int                          `json:"native_tokens_prompt,omitempty"`
	NativeCompletionTokens      int                          `json:"native_tokens_completion,omitempty"`
	NativeCompletionImageTokens int                          `json:"native_tokens_completion_images,omitempty"`
	NativeReasoningTokens       int                          `json:"native_tokens_reasoning,omitempty"`
	NativeCachedTokens          int                          `json:"native_tokens_cached,omitempty"`
	NumMediaPrompt              int                          `json:"num_media_prompt,omitempty"`
	NumInputAudioPrompt         int                          `json:"num_input_audio_prompt,omitempty"`
	NumMediaCompletion          int                          `json:"num_media_completion,omitempty"`
	NumSearchResults            int                          `json:"num_search_results,omitempty"`
	Origin                      string                       `json:"origin,omitempty"`
	Usage                       float64                      `json:"usage,omitempty"`
	IsBYOK                      bool                         `json:"is_byok,omitempty"`
	ExternalUser                string                       `json:"external_user,omitempty"`
	APIType                     string                       `json:"api_type,omitempty"`
	Router                      string                       `json:"router,omitempty"`
	ProviderResponses           []GenerationProviderResponse `json:"provider_responses,omitempty"`
}

type GenerationProviderResponse struct {
	ID             string   `json:"id,omitempty"`
	EndpointID     string   `json:"endpoint_id,omitempty"`
	ModelPermaslug string   `json:"model_permaslug,omitempty"`
	ProviderName   string   `json:"provider_name,omitempty"`
	Status         *int     `json:"status,omitempty"`
	Latency        *float64 `json:"latency,omitempty"`
	IsBYOK         *bool    `json:"is_byok,omitempty"`
}

func (c *Client) GetGeneration(ctx context.Context, id string) (*Generation, error) {
	query := url.Values{}
	query.Set("id", id)

	var res GenerationResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/generation", query, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}
