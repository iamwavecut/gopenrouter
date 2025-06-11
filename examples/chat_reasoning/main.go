package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	req := gopenrouter.ChatCompletionRequest{
		Model: "anthropic/claude-3-sonnet", // Use a model that supports reasoning
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "Think step by step. What is the capital of France?",
			},
		},
		Reasoning: &gopenrouter.ReasoningParams{
			MaxTokens: 1000,
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println("Assistant's Answer:")
	fmt.Println(resp.Choices[0].Message.Content)

	if resp.Choices[0].Message.Reasoning != "" {
		fmt.Println("\nReasoning:")
		fmt.Println(resp.Choices[0].Message.Reasoning)
	}
}
