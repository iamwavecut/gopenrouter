package gopenrouter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

func TestNewRequest(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	testCases := []struct {
		name            string
		clientConfig    ClientConfig
		method          string
		url             string
		payload         any
		expectedError   *RequestError
		expectedHeaders http.Header
		expectedBody    string
	}{
		{
			name:         "GET request without payload",
			clientConfig: ClientConfig{AuthToken: "test-token"},
			method:       http.MethodGet,
			url:          "https://example.com/test",
			payload:      nil,
			expectedHeaders: http.Header{
				"Authorization": {"Bearer test-token"},
			},
			expectedBody: "",
		},
		{
			name:         "POST request with payload",
			clientConfig: ClientConfig{AuthToken: "test-token"},
			method:       http.MethodPost,
			url:          "https://example.com/test",
			payload:      payload{Name: "test-name"},
			expectedHeaders: http.Header{
				"Authorization": {"Bearer test-token"},
				"Content-Type":  {"application/json"},
			},
			expectedBody: `{"name":"test-name"}` + "\n",
		},
		{
			name: "Request with site URL and name",
			clientConfig: ClientConfig{
				AuthToken:      "test-token",
				SiteURL:        "https://my-site.com",
				SiteName:       "My Site",
				SiteCategories: []string{"productivity", "coding"},
			},
			method:  http.MethodPost,
			url:     "https://example.com/test",
			payload: nil,
			expectedHeaders: http.Header{
				"Authorization":           {"Bearer test-token"},
				"Http-Referer":            {"https://my-site.com"},
				"X-Openrouter-Title":      {"My Site"},
				"X-Title":                 {"My Site"},
				"X-Openrouter-Categories": {"productivity,coding"},
			},
			expectedBody: "",
		},
		{
			name:          "Invalid payload that cannot be marshalled",
			clientConfig:  ClientConfig{AuthToken: "test-token"},
			method:        http.MethodPost,
			url:           "https://example.com/test",
			payload:       make(chan int),
			expectedError: &RequestError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClientWithConfig(tc.clientConfig)
			req, err := client.newRequest(context.Background(), tc.method, tc.url, tc.payload)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				if _, ok := err.(*RequestError); !ok {
					t.Fatalf("expected error of type RequestError, got %T", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if req.Method != tc.method {
				t.Errorf("expected method %q, got %q", tc.method, req.Method)
			}
			if req.URL.String() != tc.url {
				t.Errorf("expected URL %q, got %q", tc.url, req.URL.String())
			}

			for key, values := range tc.expectedHeaders {
				if !reflect.DeepEqual(req.Header[key], values) {
					t.Errorf("expected header %s = %v, got %v", key, values, req.Header[key])
				}
			}

			if req.Body != nil {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				if err := req.Body.Close(); err != nil {
					t.Fatalf("close body: %v", err)
				}
				if string(bodyBytes) != tc.expectedBody {
					t.Errorf("expected body %q, got %q", tc.expectedBody, string(bodyBytes))
				}
			} else if tc.expectedBody != "" {
				t.Errorf("expected body %q but request body is nil", tc.expectedBody)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      APIError
		expected string
	}{
		{
			name:     "error with message",
			err:      APIError{Message: "test error"},
			expected: "test error",
		},
		{
			name:     "error with empty message",
			err:      APIError{Message: ""},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("expected error message %q, got %q", tc.expected, tc.err.Error())
			}
		})
	}
}

func TestRequestError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	reqErr := &RequestError{Err: originalErr}
	expected := fmt.Sprintf("request error: %v", originalErr)
	if reqErr.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, reqErr.Error())
	}
}

func TestRequestError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	reqErr := &RequestError{Err: originalErr}
	if !errors.Is(reqErr, originalErr) {
		t.Errorf("expected unwrapped error to be the original error")
	}
}
