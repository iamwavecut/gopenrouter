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
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "Explain quantum entanglement in one paragraph."},
		},
		StreamOptions: &gopenrouter.StreamOptions{IncludeUsage: true},
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Printf("Stream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\n[DONE]")
			return
		}
		if err != nil {
			fmt.Printf("stream recv error: %v\n", err)
			return
		}
		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
		if chunk.Usage != nil {
			fmt.Printf("\n[Usage] prompt=%d completion=%d total=%d cost=$%.4f\n",
				chunk.Usage.PromptTokens, chunk.Usage.CompletionTokens, chunk.Usage.TotalTokens, chunk.Usage.Cost)
		}
	}
}
