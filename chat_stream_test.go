package gopenrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateChatCompletionStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		var reqBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if stream, ok := reqBody["stream"]; !ok || !stream.(bool) {
			t.Errorf("expected stream to be true in request body, but it was not")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {}\n\n"))
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	req := ChatCompletionRequest{
		Model: "test-model",
		Messages: []ChatCompletionMessage{
			{Role: RoleUser, Content: "request content"},
		},
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream() err = %v, want nil", err)
	}
	defer stream.Close()
}

func TestCreateChatCompletionStream_APIError(t *testing.T) {
	errorMsg := "invalid request"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, errorMsg)
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	req := ChatCompletionRequest{Model: "test-model"}
	_, err := client.CreateChatCompletionStream(context.Background(), req)

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

func TestChatCompletionStream_Recv(t *testing.T) {
	mockResponses := []ChatCompletionStreamResponse{
		{ID: "1", Choices: []ChatCompletionStreamChoice{{Delta: ChatCompletionMessage{Content: "Hello"}}}},
		{ID: "2", Choices: []ChatCompletionStreamChoice{{Delta: ChatCompletionMessage{Content: " World"}}}},
	}

	var sb strings.Builder
	for _, res := range mockResponses {
		b, _ := json.Marshal(res)
		sb.WriteString(fmt.Sprintf("data: %s\n\n", b))
	}
	sb.WriteString("data: [DONE]\n\n")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		io.WriteString(w, sb.String())
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{Model: "test"})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	for i := 0; i < len(mockResponses); i++ {
		resp, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv() error: %v", err)
		}
		if resp.ID != mockResponses[i].ID {
			t.Errorf("Expected ID %s, got %s", mockResponses[i].ID, resp.ID)
		}
	}

	_, err = stream.Recv()
	if err != io.EOF {
		t.Errorf("Expected io.EOF, got %v", err)
	}
}

func TestCreateChatCompletionStream_IncludeUsageRequested(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]any
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		// ensure stream_options.include_usage is present when set by caller
		so, ok := reqBody["stream_options"].(map[string]any)
		if !ok || so["include_usage"] != true {
			t.Errorf("expected stream_options.include_usage=true, got: %v", reqBody["stream_options"])
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {}\n\n"))
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	req := ChatCompletionRequest{Model: "test", StreamOptions: &StreamOptions{IncludeUsage: true}}
	_, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStreamLogprobsDelta_Unmarshal(t *testing.T) {
	chunk := ChatCompletionStreamResponse{
		ID: "1",
		Choices: []ChatCompletionStreamChoice{
			{
				Index: 0,
				Logprobs: &ChatCompletionStreamChoiceLogprobs{
					Content: []ChatCompletionStreamChoiceDelta{{Token: "Hello", LogProb: -0.1}},
				},
			},
		},
	}
	b, _ := json.Marshal(chunk)
	var got ChatCompletionStreamResponse
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.Choices[0].Logprobs == nil || len(got.Choices[0].Logprobs.Content) != 1 {
		t.Fatalf("expected one logprobs content delta, got: %+v", got.Choices[0].Logprobs)
	}
}
