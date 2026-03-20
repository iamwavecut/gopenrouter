package embeddings

import (
	"encoding/json"

	"github.com/iamwavecut/gopenrouter/catalog"
	"github.com/iamwavecut/gopenrouter/shared"
)

type (
	Provider            = shared.Provider
	ProviderPreferences = shared.ProviderPreferences
	ImageURL            = shared.ImageURL
	ModelsList          = catalog.ModelsList
)

type Request struct {
	Input          any                  `json:"input"`
	Model          string               `json:"model"`
	EncodingFormat string               `json:"encoding_format,omitempty"`
	Dimensions     int                  `json:"dimensions,omitempty"`
	User           string               `json:"user,omitempty"`
	Provider       *ProviderPreferences `json:"provider,omitempty"`
	InputType      string               `json:"input_type,omitempty"`
}

type InputPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type MultimodalInput struct {
	Content []InputPart `json:"content"`
}

type Response struct {
	ID     string  `json:"id,omitempty"`
	Object string  `json:"object"`
	Data   []Datum `json:"data"`
	Model  string  `json:"model"`
	Usage  *Usage  `json:"usage,omitempty"`
}

type Datum struct {
	Object    string `json:"object"`
	Embedding Value  `json:"embedding"`
	Index     int    `json:"index,omitempty"`
}

type Value struct {
	Vector []float64
	Base64 string
}

func (e Value) MarshalJSON() ([]byte, error) {
	if e.Base64 != "" {
		return json.Marshal(e.Base64)
	}
	return json.Marshal(e.Vector)
}

func (e *Value) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return json.Unmarshal(data, &e.Base64)
	}
	return json.Unmarshal(data, &e.Vector)
}

type Usage struct {
	PromptTokens int     `json:"prompt_tokens"`
	TotalTokens  int     `json:"total_tokens"`
	Cost         float64 `json:"cost,omitempty"`
}
