package gopenrouter

import (
	"net/http"
)

const (
	openRouterAPIV1 = "https://openrouter.ai/api/v1"
)

// Client is a client for the OpenRouter API.
type Client struct {
	config ClientConfig
}

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

// NewClient creates a new OpenRouter API client.
func NewClient(authToken string) *Client {
	config := DefaultConfig(authToken)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates a new OpenRouter API client with a custom configuration.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// DefaultConfig returns a default configuration for the OpenRouter API client.
func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		HTTPClient: &http.Client{},
		BaseURL:    openRouterAPIV1,
		AuthToken:  authToken,
	}
}

func (ClientConfig) String() string {
	return "<OpenRouter API ClientConfig>"
}
