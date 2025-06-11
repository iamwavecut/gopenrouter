package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCheckCredits(t *testing.T) {
	mockKeyData := KeyData{
		Label: "test-key",
		Limit: 100.0,
	}
	mockResponse := KeyCheckResponse{Data: mockKeyData}
	mockRespBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/key" {
			t.Errorf("expected path /auth/key, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockRespBody)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	keyData, err := client.CheckCredits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*keyData, mockKeyData) {
		t.Errorf("response mismatch:\n- got: %+v\n- want: %+v", *keyData, mockKeyData)
	}
}

func TestCheckCredits_APIError(t *testing.T) {
	errorMsg := "invalid key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, errorMsg)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	_, err := client.CheckCredits(context.Background())
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected error of type APIError, got %T", err)
	}

	if apiErr.Message != errorMsg {
		t.Errorf("expected error message %q, got %q", errorMsg, apiErr.Message)
	}
}
