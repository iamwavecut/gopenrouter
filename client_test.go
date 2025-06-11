package gopenrouter

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		name      string
		authToken string
	}{
		{
			name:      "empty token",
			authToken: "",
		},
		{
			name:      "non-empty token",
			authToken: "test-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.authToken)

			if client == nil {
				t.Fatal("NewClient returned nil")
			}
			if client.config.AuthToken != tc.authToken {
				t.Errorf("expected auth token %q, got %q", tc.authToken, client.config.AuthToken)
			}
			if client.config.BaseURL != openRouterAPIV1 {
				t.Errorf("expected base URL %q, got %q", openRouterAPIV1, client.config.BaseURL)
			}
			if client.config.HTTPClient == nil {
				t.Error("expected http client to be set, but it was nil")
			}
		})
	}
}

func TestNewClientWithConfig(t *testing.T) {
	customClient := &http.Client{}
	config := ClientConfig{
		AuthToken:  "custom-token",
		BaseURL:    "https://custom.url/api",
		HTTPClient: customClient,
	}

	client := NewClientWithConfig(config)

	if client == nil {
		t.Fatal("NewClientWithConfig returned nil")
	}

	if client.config.AuthToken != config.AuthToken {
		t.Errorf("expected auth token %q, got %q", config.AuthToken, client.config.AuthToken)
	}
	if client.config.BaseURL != config.BaseURL {
		t.Errorf("expected base URL %q, got %q", config.BaseURL, client.config.BaseURL)
	}
	if client.config.HTTPClient != customClient {
		t.Error("http client not set correctly")
	}
}

func TestDefaultConfig(t *testing.T) {
	authToken := "test-auth-token"
	config := DefaultConfig(authToken)

	if config.AuthToken != authToken {
		t.Errorf("expected auth token %q, got %q", authToken, config.AuthToken)
	}

	if config.BaseURL != openRouterAPIV1 {
		t.Errorf("expected base URL %q, got %q", openRouterAPIV1, config.BaseURL)
	}

	if config.HTTPClient == nil {
		t.Errorf("expected HTTPClient to be non-nil")
	}
}
