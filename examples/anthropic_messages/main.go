package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
	anthropicapi "github.com/iamwavecut/gopenrouter/anthropic"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	api := anthropicapi.New(client)

	resp, err := api.Create(context.Background(), anthropicapi.Request{
		Model:     "anthropic/claude-3.5-sonnet",
		MaxTokens: 128,
		Messages: []anthropicapi.Message{
			{
				Role:    "user",
				Content: "Explain in one sentence why provider routing matters.",
			},
		},
	})
	if err != nil {
		fmt.Printf("anthropic.Create error: %v\n", err)
		return
	}

	if len(resp.Content) == 0 {
		fmt.Println("No content blocks returned")
		return
	}
	fmt.Println(resp.Content[0].Text)
}
