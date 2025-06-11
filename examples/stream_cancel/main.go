package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	// Create a context that will be canceled after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := gopenrouter.ChatCompletionRequest{
		Model: "mistralai/mistral-7b-instruct",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "Write a very long story about a space explorer who finds a new galaxy.",
			},
		},
	}

	fmt.Println("Starting stream, will cancel in 3 seconds...")

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			// When the context is canceled, the underlying HTTP client will return an error.
			// We check for that error to confirm cancellation.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("\nStream successfully canceled.")
			} else if !errors.Is(err, io.EOF) {
				fmt.Printf("\nStream error: %v\n", err)
			} else {
				fmt.Println("\nStream finished naturally.")
			}
			return
		}

		fmt.Print(response.Choices[0].Delta.Content)
	}
}
