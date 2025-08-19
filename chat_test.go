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

func TestChatCompletionRequest_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		req      ChatCompletionRequest
		expected string
	}{
		{
			name: "with extra body",
			req: ChatCompletionRequest{
				Model: "test-model",
				ExtraBody: map[string]any{
					"extra_param": "extra_value",
				},
			},
			expected: `{"model":"test-model","extra_param":"extra_value"}`,
		},
		{
			name: "without extra body",
			req: ChatCompletionRequest{
				Model: "test-model",
			},
			expected: `{"model":"test-model"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Unmarshal and remarshal to handle key order differences
			var actualMap map[string]any
			if err := json.Unmarshal(b, &actualMap); err != nil {
				t.Fatalf("failed to unmarshal actual json: %v", err)
			}
			actualRemarshaled, _ := json.Marshal(actualMap)

			var expectedMap map[string]any
			if err := json.Unmarshal([]byte(tc.expected), &expectedMap); err != nil {
				t.Fatalf("failed to unmarshal expected json: %v", err)
			}
			expectedRemarshaled, _ := json.Marshal(expectedMap)

			if string(actualRemarshaled) != string(expectedRemarshaled) {
				t.Errorf("expected json %s, got %s", tc.expected, string(b))
			}
		})
	}
}

func TestChatCompletionMessage_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *ChatCompletionMessage
		expected string
	}{
		{
			name: "with multi-content",
			msg: &ChatCompletionMessage{
				Role: RoleUser,
				MultiContent: []ChatCompletionMessagePart{
					{Type: "text", Text: "hello"},
				},
			},
			expected: `{"role":"user","content":[{"type":"text","text":"hello"}]}`,
		},
		{
			name: "with single content",
			msg: &ChatCompletionMessage{
				Role:    RoleUser,
				Content: "hello",
			},
			expected: `{"role":"user","content":"hello"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var actualMap, expectedMap map[string]any
			if err := json.Unmarshal(b, &actualMap); err != nil {
				t.Fatalf("failed to unmarshal actual json: %v", err)
			}
			if err := json.Unmarshal([]byte(tc.expected), &expectedMap); err != nil {
				t.Fatalf("failed to unmarshal expected json: %v", err)
			}

			if !reflect.DeepEqual(actualMap, expectedMap) {
				t.Errorf("json mismatch:\ngot:    %s\nwanted: %s", string(b), tc.expected)
			}
		})
	}
}

func TestCreateChatCompletion(t *testing.T) {
	mockResponse := ChatCompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: 12345,
		Model:   "test-model",
		Choices: []Choice{
			{
				Index: 0,
				Message: ChatCompletionMessage{
					Role:    RoleAssistant,
					Content: "response content",
				},
				FinishReason: "stop",
			},
		},
	}
	mockRespBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected path /chat/completions, got %s", r.URL.Path)
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

	req := ChatCompletionRequest{
		Model: "test-model",
		Messages: []ChatCompletionMessage{
			{Role: RoleUser, Content: "request content"},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*resp, mockResponse) {
		t.Errorf("response mismatch:\n- got: %+v\n- want: %+v", *resp, mockResponse)
	}
}

func TestCreateChatCompletion_APIError(t *testing.T) {
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
	_, err := client.CreateChatCompletion(context.Background(), req)

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

func TestChatCompletionMessage_ToolCallsMarshal(t *testing.T) {
	msg := &ChatCompletionMessage{
		Role: RoleAssistant,
		ToolCalls: []ToolCall{
			{
				ID:   "call_abc123",
				Type: "function",
				Function: Function{
					Name:      "search_gutenberg_books",
					Arguments: `{"search_terms":["James","Joyce"]}`,
				},
			},
		},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	tc, ok := m["tool_calls"].([]any)
	if !ok || len(tc) != 1 {
		t.Fatalf("expected one tool_call, got: %v", m["tool_calls"])
	}
	call := tc[0].(map[string]any)
	if call["id"] != "call_abc123" {
		t.Errorf("unexpected tool_call.id: %v", call["id"])
	}
	if call["type"] != "function" {
		t.Errorf("unexpected tool_call.type: %v", call["type"])
	}
	fn := call["function"].(map[string]any)
	if fn["name"] != "search_gutenberg_books" {
		t.Errorf("unexpected function.name: %v", fn["name"])
	}
	if fn["arguments"] != `{"search_terms":["James","Joyce"]}` {
		t.Errorf("unexpected function.arguments: %v", fn["arguments"])
	}
}

func TestChatCompletionRequest_ParallelToolCalls_Marshal(t *testing.T) {
	parallel := false
	req := ChatCompletionRequest{
		Model:             "test-model",
		ParallelToolCalls: &parallel,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if v, ok := m["parallel_tool_calls"].(bool); !ok || v != false {
		t.Errorf("expected parallel_tool_calls=false, got: %v", m["parallel_tool_calls"])
	}
}

func TestChatCompletionRequest_Logprobs_Marshal(t *testing.T) {
	lp := true
	top := 5
	req := ChatCompletionRequest{
		Model:       "test-model",
		LogProbs:    &lp,
		TopLogProbs: &top,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if v, ok := m["logprobs"].(bool); !ok || v != true {
		t.Errorf("expected logprobs=true, got: %v", m["logprobs"])
	}
	if v, ok := m["top_logprobs"].(float64); !ok || int(v) != 5 {
		t.Errorf("expected top_logprobs=5, got: %v", m["top_logprobs"])
	}
}

func TestCreateChatCompletion_NonJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "oops plain text")
	}))
	defer server.Close()

	config := DefaultConfig("test-token")
	config.BaseURL = server.URL
	client := NewClientWithConfig(config)

	_, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected RequestError, got %T", err)
	}
	if reqErr.HTTPStatusCode != http.StatusInternalServerError {
		t.Errorf("expected HTTPStatusCode=%d, got %d", http.StatusInternalServerError, reqErr.HTTPStatusCode)
	}
	if reqErr.HTTPStatus == "" {
		t.Errorf("expected HTTPStatus to be set")
	}
}
