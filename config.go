package gopenrouter

import "net/http"

const (
	openRouterAPIV1 = "https://openrouter.ai/api/v1"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	AuthToken string
	BaseURL   string
	// Deprecated: this field is kept only for backward compatibility and is not sent, because the current public OpenRouter API does not define an organization header.
	OrgID          string
	HTTPClient     *http.Client
	SiteURL        string
	SiteName       string
	SiteCategories []string

	// Deprecated: use AuthToken instead.
	APIKey string
}

func (ClientConfig) String() string {
	return "<OpenRouter API ClientConfig>"
}
