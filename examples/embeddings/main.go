package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
	embeddingsapi "github.com/iamwavecut/gopenrouter/embeddings"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	api := embeddingsapi.New(client)

	resp, err := api.Create(context.Background(), embeddingsapi.Request{
		Model: "openai/text-embedding-3-small",
		Input: []string{
			"OpenRouter embeddings example",
			"gopenrouter SDK sync",
		},
	})
	if err != nil {
		fmt.Printf("embeddings.Create error: %v\n", err)
		return
	}

	fmt.Printf("Embeddings: %d\n", len(resp.Data))
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Prompt tokens: %d\n", resp.Usage.PromptTokens)
}
