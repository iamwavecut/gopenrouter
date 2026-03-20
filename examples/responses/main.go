package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
	responsesapi "github.com/iamwavecut/gopenrouter/responses"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	api := responsesapi.New(client)

	resp, err := api.Create(context.Background(), responsesapi.Request{
		Model: "openai/gpt-4o-mini",
		Input: []map[string]any{
			{
				"type": "message",
				"role": "user",
				"content": []map[string]any{
					{
						"type": "input_text",
						"text": "Give me a one-line summary of OpenRouter.",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("responses.Create error: %v\n", err)
		return
	}

	fmt.Printf("Response ID: %s\n", resp.ID)
	fmt.Printf("Status: %s\n", resp.Status)
	if len(resp.Output) > 0 && len(resp.Output[0].Content) > 0 {
		fmt.Println(resp.Output[0].Content[0].Text)
	}
}
