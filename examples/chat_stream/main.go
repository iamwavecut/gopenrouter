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
		Model: "gryphe/mythomax-l2-13b",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello, what is the capital of France?",
			},
		},
	}

	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
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

		fmt.Print(response.Choices[0].Delta.Content)
	}
}
