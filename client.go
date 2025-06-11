package gopenrouter

import (
	"net/http"
)

// Client is a client for the OpenRouter API.
type Client struct {
	config ClientConfig
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
