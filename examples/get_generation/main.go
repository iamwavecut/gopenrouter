package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	ctx := context.Background()

	// First, create a generation to get an ID
	chatReq := gopenrouter.ChatCompletionRequest{
		Model: "google/gemini-2.0-flash-lite-001",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "What is the airspeed velocity of an unladen swallow?",
			},
		},
	}

	chatResp, err := client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	if chatResp.ID == "" {
		fmt.Println("Failed to get a generation ID from chat completion.")
		return
	}

	fmt.Printf("Got generation ID: %s\n\n", chatResp.ID)

	// Now, use the ID to get generation stats
	// It may take a moment for the generation to be available, so we'll retry a few times.
	var genResp *gopenrouter.Generation
	var getErr error
	for i := 0; i < 5; i++ {
		genResp, getErr = client.GetGeneration(ctx, chatResp.ID)
		if getErr == nil {
			break
		}
		fmt.Printf("Attempt %d: Failed to get generation, retrying in 1 second... (error: %v)\n", i+1, getErr)
		time.Sleep(1 * time.Second)
	}

	if getErr != nil {
		fmt.Printf("GetGeneration error after retries: %v\n", getErr)
		return
	}

	if genResp == nil {
		fmt.Println("Failed to retrieve generation data, response was nil.")
		return
	}

	fmt.Printf("Successfully retrieved generation:\n")
	fmt.Printf("  ID: %s\n", genResp.ID)
	fmt.Printf("  Model: %s\n", genResp.Model)
	fmt.Printf("  Prompt Tokens: %d\n", genResp.PromptTokens)
	fmt.Printf("  Completion Tokens: %d\n", genResp.CompletionTokens)
	fmt.Printf("  Total Tokens: %d\n", genResp.PromptTokens+genResp.CompletionTokens)
	fmt.Printf("  Cost: %f\n", genResp.TotalCost)
}
