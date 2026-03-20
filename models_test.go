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
		if _, err := w.Write(mockRespBody); err != nil {
			t.Fatalf("write response: %v", err)
		}
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

func TestPlugins_Marshal(t *testing.T) {
	req := ChatCompletionRequest{
		Model: "test-model",
		Plugins: []Plugin{
			{ID: PluginIDWeb, Config: WebSearchOptions{SearchContextSize: SearchContextSizeHigh}},
			{ID: PluginIDFileParser, Config: FileParserConfig{PDF: &PDFPlugin{Engine: string(PDFEnginePDFText)}}},
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	plugins, ok := m["plugins"].([]any)
	if !ok || len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got: %v", m["plugins"])
	}
	webPlugin := plugins[0].(map[string]any)
	if _, ok := webPlugin["config"]; ok {
		t.Fatalf("expected flattened plugin payload, got nested config: %v", webPlugin)
	}
	if webPlugin["search_context_size"] != string(SearchContextSizeHigh) {
		t.Fatalf("expected legacy config bridge to preserve search_context_size, got %v", webPlugin["search_context_size"])
	}
	filePlugin := plugins[1].(map[string]any)
	if _, ok := filePlugin["pdf"]; !ok {
		t.Fatalf("expected file-parser pdf config at top level, got %v", filePlugin)
	}
}
