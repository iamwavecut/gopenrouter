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

	// Prepare a tiny PDF payload (dummy content)
	fileContent := []byte("%PDF-1.1\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF\n")
	encodedFile := base64.StdEncoding.EncodeToString(fileContent)

	req := gopenrouter.ChatCompletionRequest{
		Model: "anthropic/claude-4-sonnet",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role: gopenrouter.RoleUser,
				MultiContent: []gopenrouter.ChatCompletionMessagePart{
					{Type: "text", Text: "Summarize the content of the attached PDF."},
					{Type: "file", File: &gopenrouter.File{Filename: "sample.pdf", FileData: "data:application/pdf;base64," + encodedFile}},
				},
			},
		},
		Plugins: []gopenrouter.Plugin{{
			ID:     gopenrouter.PluginIDFileParser,
			Config: gopenrouter.FileParserConfig{PDF: &gopenrouter.PDFPlugin{Engine: string(gopenrouter.PDFEnginePDFText)}},
		}},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
