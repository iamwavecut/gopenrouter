package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

// A tiny 1x1 red pixel PNG as a base64 string.
const redPixelBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAwAB/epv2AAAAABJRU5ErkJggg=="

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	req := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o-mini", // Use a model that supports vision
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role: gopenrouter.RoleUser,
				MultiContent: []gopenrouter.ChatCompletionMessagePart{
					{
						Type: "text",
						Text: "What color is this image?",
					},
					{
						Type: "image_url",
						ImageURL: &gopenrouter.ImageURL{
							URL: "data:image/png;base64," + redPixelBase64,
						},
					},
				},
			},
		},
		MaxTokens: 10,
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
