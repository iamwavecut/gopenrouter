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
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "Find the latest news headlines about Go programming."},
		},
		Plugins: []gopenrouter.Plugin{{
			ID:     gopenrouter.PluginIDWeb,
			Config: gopenrouter.WebSearchOptions{SearchContextSize: gopenrouter.SearchContextSizeHigh},
		}},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
