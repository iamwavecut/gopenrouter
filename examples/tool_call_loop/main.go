package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	tools := []gopenrouter.Tool{
		{
			Type: "function",
			Function: gopenrouter.Function{
				Name:        "search_gutenberg_books",
				Description: "Search for books in the Project Gutenberg library",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"search_terms": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "List of search terms to find books",
						},
					},
					"required": []string{"search_terms"},
				},
			},
		},
	}

	req1 := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "What are the titles of some James Joyce books?"},
		},
		Tools: tools,
		ToolChoice: map[string]any{
			"type":     "function",
			"function": map[string]any{"name": "search_gutenberg_books"},
		},
	}

	ctx := context.Background()
	resp1, err := client.CreateChatCompletion(ctx, req1)
	if err != nil {
		fmt.Printf("ChatCompletion step 1 error: %v\n", err)
		return
	}

	if len(resp1.Choices) == 0 || len(resp1.Choices[0].Message.ToolCalls) == 0 {
		fmt.Println(resp1.Choices[0].Message.Content)
		return
	}

	toolCall := resp1.Choices[0].Message.ToolCalls[0]

	// Simulate executing the tool locally
	// Optionally parse arguments for real use: json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	toolResult := []map[string]any{
		{"id": 4300, "title": "Ulysses", "authors": []map[string]string{{"name": "Joyce, James"}}},
		{"id": 2814, "title": "Dubliners", "authors": []map[string]string{{"name": "Joyce, James"}}},
	}
	toolResultJSON, _ := json.Marshal(toolResult)

	req2 := gopenrouter.ChatCompletionRequest{
		Model: req1.Model,
		Tools: tools,
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "What are the titles of some James Joyce books?"},
			{Role: gopenrouter.RoleAssistant, ToolCalls: []gopenrouter.ToolCall{toolCall}},
			{Role: gopenrouter.RoleTool, ToolCallID: toolCall.ID, Content: string(toolResultJSON)},
		},
	}

	resp2, err := client.CreateChatCompletion(ctx, req2)
	if err != nil {
		fmt.Printf("ChatCompletion step 3 error: %v\n", err)
		return
	}

	if len(resp2.Choices) > 0 {
		fmt.Println(resp2.Choices[0].Message.Content)
	}
}
