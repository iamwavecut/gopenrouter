package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
	"github.com/iamwavecut/gopenrouter/shared"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	req := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "Find the latest news headlines about Go programming."},
		},
		Plugins: []shared.Plugin{{
			ID:     shared.PluginIDWeb,
			Config: shared.WebSearchOptions{SearchContextSize: shared.SearchContextSizeHigh},
		}},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
