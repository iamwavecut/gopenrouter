package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	req := gopenrouter.ChatCompletionRequest{
		Model: "deepseek/deepseek-r1-0528",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "Please write a short story about a robot who learns to paint.",
			},
		},
		Reasoning: &gopenrouter.ReasoningParams{
			Exclude: false,
			Effort:  "low", // Request thinking tokens.
		},
	}

	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		if response.Choices[0].Delta.Reasoning != "" {
			fmt.Printf("[Thinking]: %s", response.Choices[0].Delta.Reasoning)
		}

		if response.Choices[0].Delta.Content != "" {
			fmt.Printf("%s", response.Choices[0].Delta.Content)
		}
	}
}
