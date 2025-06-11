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

	hugeTextBody := strings.Repeat("This is a large block of text that we want to cache. ", 500)

	req := gopenrouter.ChatCompletionRequest{
		Model: "anthropic/claude-3.5-sonnet",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role: gopenrouter.RoleUser,
				MultiContent: []gopenrouter.ChatCompletionMessagePart{
					{
						Type: "text",
						Text: "Given the book below:",
					},
					{
						Type: "text",
						Text: hugeTextBody,
						CacheControl: &gopenrouter.CacheControl{
							Type: "ephemeral",
						},
					},
					{
						Type: "text",
						Text: "Name all the characters in the above book",
					},
				},
			},
		},
		MaxTokens: 100,
		Usage: &gopenrouter.UsageParams{
			Include: true,
		},
	}

	fmt.Println("Sending request with cache control...")

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Chat completion error: %v\n", err)
	}

	fmt.Printf("Response from model: %s\n", resp.Choices[0].Message.Content)
	if resp.Usage.PromptTokensDetails != nil {
		fmt.Printf("Cached tokens: %d\n", resp.Usage.PromptTokensDetails.CachedTokens)
	}
	fmt.Printf("Total cost: %f\n", resp.Usage.Cost)

	// Now, send a follow-up request. If caching worked, the cost should be lower.
	// Note: The cache is ephemeral and might expire quickly. This is for demonstration.

	fmt.Println("\nSending a follow-up request to test cache read...")

	followUpReq := gopenrouter.ChatCompletionRequest{
		Model: "anthropic/claude-3.5-sonnet",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role: gopenrouter.RoleUser,
				MultiContent: []gopenrouter.ChatCompletionMessagePart{
					{
						Type: "text",
						Text: "Given the book below:",
					},
					{
						Type: "text",
						Text: hugeTextBody,
						CacheControl: &gopenrouter.CacheControl{
							Type: "ephemeral",
						},
					},
					{
						Type: "text",
						Text: "What is the main theme of the book described?",
					},
				},
			},
		},
		MaxTokens: 100,
		Usage: &gopenrouter.UsageParams{
			Include: true,
		},
	}

	followUpResp, err := client.CreateChatCompletion(context.Background(), followUpReq)
	if err != nil {
		log.Fatalf("Follow-up chat completion error: %v\n", err)
	}

	fmt.Printf("Follow-up response from model: %s\n", followUpResp.Choices[0].Message.Content)
	if followUpResp.Usage.PromptTokensDetails != nil {
		fmt.Printf("Cached tokens on follow-up: %d\n", followUpResp.Usage.PromptTokensDetails.CachedTokens)
	}
	fmt.Printf("Follow-up total cost: %f\n", followUpResp.Usage.Cost)

	fmt.Println("\nNote: Cost savings depend on the provider's caching policy and pricing.")
	fmt.Println("Anthropic charges for cache writes, but reads are cheaper.")
	fmt.Println("Check your OpenRouter dashboard for detailed cost breakdowns.")
}
