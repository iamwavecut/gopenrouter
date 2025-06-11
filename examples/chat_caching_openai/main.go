package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	// OpenAI prompt caching is automatic for prompts over 1024 tokens.
	// No special `cache_control` field is needed.
	largePrompt := strings.Repeat("This is a large block of text to demonstrate OpenAI's automatic prompt caching. ", 200)

	req := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o-mini", // Or another supported OpenAI model
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    gopenrouter.RoleUser,
				Content: largePrompt + "\n\nQuestion: What is the purpose of the text above?",
			},
		},
		MaxTokens: 100,
		Usage: &gopenrouter.UsageParams{
			Include: true,
		},
	}

	fmt.Println("Sending first request to OpenAI...")
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Chat completion error: %v\n", err)
	}

	fmt.Printf("Response from model: %s\n", resp.Choices[0].Message.Content)
	if resp.Usage.PromptTokensDetails != nil {
		fmt.Printf("Cached tokens: %d\n", resp.Usage.PromptTokensDetails.CachedTokens)
	}
	fmt.Printf("Total cost: %f\n", resp.Usage.Cost)

	// Send the exact same request again. The second time should be cheaper due to caching.
	fmt.Println("\nSending second request to OpenAI to test cache read...")
	followUpResp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Follow-up chat completion error: %v\n", err)
	}

	fmt.Printf("Follow-up response from model: %s\n", followUpResp.Choices[0].Message.Content)
	if followUpResp.Usage.PromptTokensDetails != nil {
		fmt.Printf("Cached tokens on follow-up: %d\n", followUpResp.Usage.PromptTokensDetails.CachedTokens)
	}
	fmt.Printf("Follow-up total cost: %f\n", followUpResp.Usage.Cost)

	fmt.Println("\nNote: For OpenAI, caching is automatic for large prompts.")
	fmt.Println("The second request should show a lower cost if the prompt was cached successfully.")
	fmt.Println("Check your OpenRouter dashboard for detailed cost breakdowns.")
}
