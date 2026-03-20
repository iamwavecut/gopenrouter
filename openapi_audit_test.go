package gopenrouter

import (
	"encoding/json"
	"os"
	"testing"
)

func TestOpenAPISnapshotContainsTrackedContracts(t *testing.T) {
	data, err := os.ReadFile("testdata/openrouter_openapi_2026-03-20.json")
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}

	var doc struct {
		Paths      map[string]json.RawMessage `json:"paths"`
		Components struct {
			Schemas map[string]json.RawMessage `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("decode snapshot: %v", err)
	}

	requiredPaths := []string{
		"/chat/completions",
		"/generation",
		"/key",
		"/embeddings",
		"/responses",
		"/messages",
		"/providers",
		"/models",
		"/models/count",
		"/models/user",
		"/models/{author}/{slug}/endpoints",
		"/endpoints/zdr",
		"/credits",
		"/credits/coinbase",
		"/activity",
		"/keys",
		"/guardrails",
	}
	for _, path := range requiredPaths {
		if _, ok := doc.Paths[path]; !ok {
			t.Fatalf("snapshot missing required path %s", path)
		}
	}

	requiredSchemas := []string{
		"ChatGenerationParams",
		"ProviderPreferences",
		"Message",
		"ChatStreamingResponseChunk",
		"ModelsListResponse",
		"AnthropicMessagesRequest",
		"AnthropicMessagesResponse",
		"OpenResponsesRequest",
		"OpenResponsesNonStreamingResponse",
	}
	for _, schema := range requiredSchemas {
		if _, ok := doc.Components.Schemas[schema]; !ok {
			t.Fatalf("snapshot missing required schema %s", schema)
		}
	}
}
