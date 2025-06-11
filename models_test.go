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

func TestListModels(t *testing.T) {
	mockModelsList := ModelsList{
		Data: []Model{
			{ID: "model1", Name: "Test Model 1"},
			{ID: "model2", Name: "Test Model 2"},
		},
	}
	mockRespBody, _ := json.Marshal(mockModelsList)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Errorf("expected path /models, got %s", r.URL.Path)
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

	modelsList, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*modelsList, mockModelsList) {
		t.Errorf("response mismatch:\n- got: %+v\n- want: %+v", *modelsList, mockModelsList)
	}
}

func TestListModels_APIError(t *testing.T) {
	errorMsg := "internal server error"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, errorMsg)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	_, err := client.ListModels(context.Background())
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
