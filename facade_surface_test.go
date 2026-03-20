package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCreateEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embeddings" {
			t.Fatalf("expected /embeddings, got %s", r.URL.Path)
		}
		var req EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != "openai/text-embedding-3-small" {
			t.Fatalf("unexpected model: %s", req.Model)
		}
		fmt.Fprint(w, `{"object":"list","data":[{"object":"embedding","embedding":[0.1,0.2],"index":0}],"model":"openai/text-embedding-3-small","usage":{"prompt_tokens":3,"total_tokens":3}}`)
	}))
	defer server.Close()

	cfg := DefaultConfig("test-token")
	cfg.BaseURL = server.URL
	client := NewClientWithConfig(cfg)

	resp, err := client.CreateEmbeddings(context.Background(), EmbeddingRequest{
		Model: "openai/text-embedding-3-small",
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || len(resp.Data[0].Embedding.Vector) != 2 {
		t.Fatalf("unexpected embeddings response: %+v", resp)
	}
}

func TestResponsesAndMessagesStreams(t *testing.T) {
	responseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/responses":
			w.Header().Set("Content-Type", "text/event-stream")
			if _, err := io.WriteString(w, "data: {\"type\":\"response.created\",\"sequence_number\":0,\"response\":{\"id\":\"resp_1\",\"status\":\"in_progress\"}}\n\n"); err != nil {
				t.Fatalf("write response stream event: %v", err)
			}
			if _, err := io.WriteString(w, "data: [DONE]\n\n"); err != nil {
				t.Fatalf("write response stream done: %v", err)
			}
		case "/messages":
			w.Header().Set("Content-Type", "text/event-stream")
			if _, err := io.WriteString(w, "event: message_start\n"); err != nil {
				t.Fatalf("write messages event name: %v", err)
			}
			if _, err := io.WriteString(w, "data: {\"type\":\"message_start\"}\n\n"); err != nil {
				t.Fatalf("write messages stream event: %v", err)
			}
			if _, err := io.WriteString(w, "data: [DONE]\n\n"); err != nil {
				t.Fatalf("write messages stream done: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer responseServer.Close()

	cfg := DefaultConfig("test-token")
	cfg.BaseURL = responseServer.URL
	client := NewClientWithConfig(cfg)

	responseStream, err := client.CreateResponseStream(context.Background(), ResponseRequest{Model: "openai/gpt-4o", Input: "hello"})
	if err != nil {
		t.Fatalf("unexpected response stream error: %v", err)
	}
	defer responseStream.Close()
	responseEvent, err := responseStream.Recv()
	if err != nil {
		t.Fatalf("unexpected response event error: %v", err)
	}
	if responseEvent.Type != "response.created" || responseEvent.Response == nil || responseEvent.Response.ID != "resp_1" {
		t.Fatalf("unexpected response event: %+v", responseEvent)
	}

	messageStream, err := client.CreateAnthropicMessageStream(context.Background(), AnthropicMessageRequest{Model: "anthropic/claude-3.5-sonnet", Messages: []AnthropicMessage{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("unexpected anthropic stream error: %v", err)
	}
	defer messageStream.Close()
	messageEvent, err := messageStream.Recv()
	if err != nil {
		t.Fatalf("unexpected anthropic event error: %v", err)
	}
	if messageEvent.Event != "message_start" || messageEvent.Type != "message_start" {
		t.Fatalf("unexpected anthropic event: %+v", messageEvent)
	}
}

func TestMetadataAndManagementMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/providers":
			fmt.Fprint(w, `{"data":[{"name":"OpenAI","slug":"openai"}]}`)
		case r.URL.Path == "/models/count":
			fmt.Fprint(w, `{"data":{"count":321}}`)
		case r.URL.Path == "/credits":
			fmt.Fprint(w, `{"data":{"total_credits":100,"total_usage":25}}`)
		case r.URL.Path == "/activity":
			if r.URL.Query().Get("date") != "2025-08-24" {
				t.Fatalf("expected activity date query, got %q", r.URL.Query().Get("date"))
			}
			fmt.Fprint(w, `{"data":[{"date":"2025-08-24","model":"openai/gpt-4.1","model_permaslug":"openai/gpt-4.1-2025-04-14","endpoint_id":"ep1","provider_name":"OpenAI","usage":0.1,"byok_usage_inference":0.05,"requests":1,"prompt_tokens":10,"completion_tokens":20,"reasoning_tokens":5}]}`)
		case r.URL.Path == "/auth/keys/code":
			fmt.Fprint(w, `{"data":{"id":"auth_1","app_id":7,"created_at":"2025-08-24T10:30:00Z"}}`)
		case r.URL.Path == "/auth/keys":
			fmt.Fprint(w, `{"key":"sk-or-v1-test","user_id":"user_1"}`)
		case r.URL.Path == "/keys" && r.Method == http.MethodGet:
			fmt.Fprint(w, `{"data":[{"hash":"hash_1","name":"prod","label":"sk-or-v1-1","disabled":false}]}`)
		case r.URL.Path == "/guardrails" && r.Method == http.MethodGet:
			fmt.Fprint(w, `{"data":[{"id":"g1","name":"Default","created_at":"2025-08-24T10:30:00Z"}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := DefaultConfig("test-token")
	cfg.BaseURL = server.URL
	client := NewClientWithConfig(cfg)

	providers, err := client.ListProviders(context.Background())
	if err != nil || len(providers.Data) != 1 {
		t.Fatalf("unexpected providers result: %+v err=%v", providers, err)
	}
	count, err := client.CountModels(context.Background())
	if err != nil || count != 321 {
		t.Fatalf("unexpected models count: %d err=%v", count, err)
	}
	credits, err := client.GetCredits(context.Background())
	if err != nil || credits.TotalCredits != 100 {
		t.Fatalf("unexpected credits: %+v err=%v", credits, err)
	}
	activity, err := client.GetUserActivity(context.Background(), ActivityParams{Date: "2025-08-24"})
	if err != nil || len(activity) != 1 {
		t.Fatalf("unexpected activity: %+v err=%v", activity, err)
	}
	authCode, err := client.CreateAuthCode(context.Background(), CreateAuthCodeRequest{CallbackURL: "https://example.com/callback"})
	if err != nil || authCode.ID != "auth_1" {
		t.Fatalf("unexpected auth code: %+v err=%v", authCode, err)
	}
	apiKey, err := client.ExchangeAuthCodeForAPIKey(context.Background(), ExchangeAuthCodeRequest{Code: "auth_1"})
	if err != nil || !strings.HasPrefix(apiKey.Key, "sk-or-v1") {
		t.Fatalf("unexpected api key exchange: %+v err=%v", apiKey, err)
	}
	keys, err := client.ListAPIKeys(context.Background(), APIKeysListParams{})
	if err != nil || len(keys) != 1 {
		t.Fatalf("unexpected api keys: %+v err=%v", keys, err)
	}
	guardrails, err := client.ListGuardrails(context.Background())
	if err != nil || len(guardrails) != 1 {
		t.Fatalf("unexpected guardrails: %+v err=%v", guardrails, err)
	}
}

func TestCheckCredits(t *testing.T) {
	mockKeyData := KeyData{
		Label: "test-key",
		Limit: 100.0,
	}
	mockResponse := KeyCheckResponse{Data: mockKeyData}
	mockRespBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/key" {
			t.Errorf("expected path /key, got %s", r.URL.Path)
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

func TestCheckCredits_FallbackToLegacyPath(t *testing.T) {
	mockKeyData := KeyData{
		Label: "legacy-key",
		Limit: 100,
	}
	mockResponse := KeyCheckResponse{Data: mockKeyData}
	mockRespBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/key":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":{"message":"not found"}}`)
		case "/auth/key":
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(mockRespBody); err != nil {
				t.Fatalf("write legacy response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
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
		t.Fatalf("response mismatch: got %+v want %+v", *keyData, mockKeyData)
	}
}
