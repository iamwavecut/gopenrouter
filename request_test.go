package gopenrouter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

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
				AuthToken: "test-token",
				SiteURL:   "https://my-site.com",
				SiteName:  "My Site",
			},
			method:  http.MethodPost,
			url:     "https://example.com/test",
			payload: nil,
			expectedHeaders: http.Header{
				"Authorization": {"Bearer test-token"},
				"Http-Referer":  {"https://my-site.com"},
				"X-Title":       {"My Site"},
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
				return // Error was expected
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

			// Check headers
			for key, values := range tc.expectedHeaders {
				if !reflect.DeepEqual(req.Header[key], values) {
					t.Errorf("expected header %s = %v, got %v", key, values, req.Header[key])
				}
			}

			if req.Body != nil {
				bodyBytes, _ := ioutil.ReadAll(req.Body)
				req.Body.Close()
				if string(bodyBytes) != tc.expectedBody {
					t.Errorf("expected body %q, got %q", tc.expectedBody, string(bodyBytes))
				}
			} else if tc.expectedBody != "" {
				t.Errorf("expected body %q but request body is nil", tc.expectedBody)
			}
		})
	}
}

// Custom unmarshaling for testing to handle nil error case
func (e *RequestError) UnmarshalJSON(data []byte) error {
	var errWrapper struct {
		Err json.RawMessage `json:"err"`
	}
	if err := json.Unmarshal(data, &errWrapper); err != nil {
		return err
	}
	if string(errWrapper.Err) == "null" {
		e.Err = nil
		return nil
	}
	return json.Unmarshal(errWrapper.Err, &e.Err)
}
