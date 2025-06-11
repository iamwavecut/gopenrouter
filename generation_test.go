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

func TestGetGeneration(t *testing.T) {
	mockGeneration := Generation{
		ID:    "test-id",
		Model: "test-model",
	}
	mockResponse := GenerationResponse{Data: mockGeneration}
	mockRespBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/generation" {
			t.Errorf("expected path /generation, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("id") != "test-id" {
			t.Errorf("expected id query param to be 'test-id', got %s", r.URL.Query().Get("id"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockRespBody)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	generation, err := client.GetGeneration(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*generation, mockGeneration) {
		t.Errorf("response mismatch:\n- got: %+v\n- want: %+v", *generation, mockGeneration)
	}
}

func TestGetGeneration_APIError(t *testing.T) {
	errorMsg := "generation not found"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, errorMsg)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	_, err := client.GetGeneration(context.Background(), "test-id")
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
