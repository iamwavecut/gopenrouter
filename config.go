package gopenrouter

import "net/http"

const (
	openRouterAPIV1 = "https://openrouter.ai/api/v1"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	AuthToken  string
	BaseURL    string
	OrgID      string
	HTTPClient *http.Client
	SiteURL    string // The URL of your app or site.
	SiteName   string // The name of your app or site.

	// Deprecated: use AuthToken instead.
	APIKey string
}

func (ClientConfig) String() string {
	return "<OpenRouter API ClientConfig>"
}
