package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	// Create some dummy file content and base64 encode it
	fileContent := []byte("%PDF-1.1\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF\n")
	encodedFile := base64.StdEncoding.EncodeToString(fileContent)

	req := gopenrouter.ChatCompletionRequest{
		Model: "anthropic/claude-3.7-sonnet:beta", // Use a model that supports file attachments
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role: gopenrouter.RoleUser,
				MultiContent: []gopenrouter.ChatCompletionMessagePart{
					{
						Type: "text",
						Text: "Please summarize the content of the attached file.",
					},
					{
						Type: "file",
						File: &gopenrouter.File{
							Filename: "summary.pdf",
							FileData: "data:application/pdf;base64," + encodedFile,
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
