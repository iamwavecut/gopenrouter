# gopenrouter

[![Go Reference](https://pkg.go.dev/badge/github.com/iamwavecut/gopenrouter.svg)](https://pkg.go.dev/github.com/iamwavecut/gopenrouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/gopenrouter)](https://goreportcard.com/report/github.com/iamwavecut/gopenrouter)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-6e45e5.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

An unofficial Go client for the [OpenRouter API](https://openrouter.ai).

This library provides a comprehensive, `go-openai`-inspired interface for interacting with the OpenRouter API, giving you access to a multitude of LLMs through a single, unified client.

## Credits

This library's design and structure are heavily inspired by the excellent [go-openai](https://github.com/sashabaranov/go-openai) library.

## Key Differences from `go-openai`

While this library maintains a familiar, `go-openai`-style interface, it includes several key features and fields specifically tailored for the OpenRouter API:

- **Multi-Provider Models**: Seamlessly switch between models from different providers (e.g., Anthropic, Google, Mistral) by changing the `Model` string.
- **Cost Tracking**: The `Usage` object in responses includes a `Cost` field, providing direct access to the dollar cost of a generation.
- **Native Token Counts**: The `GetGeneration` endpoint provides access to `NativePromptTokens` and `NativeCompletionTokens`, giving you precise, provider-native tokenization data.
- **Advanced Routing**: Use `Models` for fallback chains and `Route` for custom routing logic.
- **Reasoning Parameters**: Control and request "thinking" tokens from supported models using the `Reasoning` parameters.
- **Provider-Specific `ExtraBody`**: Pass custom, provider-specific parameters through the `ExtraBody` field for fine-grained control.
- **Client Utilities**: Includes built-in methods to `ListModels`, `CheckCredits`, and `GetGeneration` stats directly from the client.

## Installation

```bash
go get github.com/iamwavecut/gopenrouter
```

## Usage

### Basic Chat Completion

```go
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
		Model: "gryphe/mythomax-l2-13b",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "Hello, what is the capital of France?",
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
```

### Streaming Chat Completion

```go
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
				Role:    gopenrouter.RoleUser,
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
```

### Vision (Image Attachments)

```go
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
```

### File Attachments

```go
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
		Model: "anthropic/claude-3.5-sonnet", // Use a model that supports file attachments
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
```

### Structured Outputs

This library supports two methods for getting structured JSON outputs from compatible models.

#### Simple JSON Mode (`json_object`)

For models that support it, you can request a JSON object without enforcing a specific schema. This is useful for simple tasks where you just need a JSON response.

```go
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
			{
				Role:    gopenrouter.RoleSystem,
				Content: "You are a helpful assistant that outputs JSON.",
			},
			{
				Role:    gopenrouter.RoleUser,
				Content: "Extract the user's name and age from this text: John Doe is 30 years old.",
			},
		},
		ResponseFormat: &gopenrouter.ResponseFormat{Type: "json_object"},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

#### Schema-Enforced JSON Mode (`json_schema`)

For more complex or critical tasks, you can provide a specific JSON schema that the model's output must conform to. The library provides a `GenerateSchema` helper to create this schema from a Go struct automatically.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

// User struct defines the desired structured output.
type User struct {
	Name string `json:"name" jsonschema:"The user's full name"`
	Age  int    `json:"age" jsonschema:"The user's age"`
}

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	// 1. Generate a JSON schema from the struct.
	schema, err := gopenrouter.GenerateSchema(User{})
	if err != nil {
		fmt.Printf("GenerateSchema error: %v\n", err)
		return
	}

	// 2. Create the request with the schema-enforced response format.
	req := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "My name is John Doe and I am 30 years old.",
			},
		},
		ResponseFormat: &gopenrouter.ResponseFormat{
			Type: "json_schema",
			JSONSchema: &gopenrouter.JSONSchema{
				Name:   "extract_user",
				Strict: true,
				Schema: schema,
			},
		},
	}

	// 3. Make the API call
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	// 4. Unmarshal the structured JSON response
	fmt.Println("Raw JSON Output:")
	fmt.Println(resp.Choices[0].Message.Content)

	var user User
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &user); err == nil {
		fmt.Printf("\nParsed User: %+v\n", user)
	} else {
		fmt.Printf("\nError parsing JSON: %v\n", err)
	}
}
```

### Advanced Usage

#### Reasoning Tokens

```go
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
		Model: "anthropic/claude-3-sonnet", // Use a model that supports reasoning
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "Think step by step. What is the capital of France?",
			},
		},
		Reasoning: &gopenrouter.ReasoningParams{
			MaxTokens: 1000,
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println("Assistant's Answer:")
	fmt.Println(resp.Choices[0].Message.Content)

	if resp.Choices[0].Message.Reasoning != "" {
		fmt.Println("\nReasoning:")
		fmt.Println(resp.Choices[0].Message.Reasoning)
	}
}
```

#### Provider-Specific Parameters with `ExtraBody`

```go
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
		Model: "google/gemini-flash-1.5",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "What is the capital of Canada?",
			},
		},
		ExtraBody: map[string]any{
			"provider": map[string]any{
				"require_parameters": true,
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
```

### Client Utilities

#### Listing Models

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	ctx := context.Background()

	models, err := client.ListModels(ctx)
	if err != nil {
		fmt.Printf("ListModels error: %v\n", err)
		return
	}

	fmt.Println("Available Models:")
	for _, model := range models.Data {
		fmt.Printf("- ID: %s, Name: %s\n", model.ID, model.Name)
	}
}
```

#### Checking Credits

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	keyData, err := client.CheckCredits(context.Background())
	if err != nil {
		fmt.Printf("CheckCredits error: %v\n", err)
		return
	}

	fmt.Println("API Key Details:")
	fmt.Printf("  Label: %s\n", keyData.Label)
	fmt.Printf("  Usage: $%f\n", keyData.Usage)
	fmt.Printf("  Limit: $%f\n", keyData.Limit)
	fmt.Printf("  Free Tier: %t\n", keyData.IsFreeTier)
}
```

#### Fetching Generation Stats

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	ctx := context.Background()

	// First, create a generation to get an ID
	chatReq := gopenrouter.ChatCompletionRequest{
		Model: "google/gemini-flash-1.5",
		Messages: []gopenrouter.ChatCompletionMessage{
			{
				Role:    gopenrouter.RoleUser,
				Content: "What is the airspeed velocity of an unladen swallow?",
			},
		},
	}

	chatResp, err := client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	if chatResp.ID == "" {
		fmt.Println("Failed to get a generation ID from chat completion.")
		return
	}

	fmt.Printf("Got generation ID: %s\n\n", chatResp.ID)

	// Now, use the ID to get generation stats
	// It may take a moment for the generation to be available, so we'll retry a few times.
	var genResp *gopenrouter.Generation
	var getErr error
	for i := 0; i < 5; i++ {
		genResp, getErr = client.GetGeneration(ctx, chatResp.ID)
		if getErr == nil {
			break
		}
		fmt.Printf("Attempt %d: Failed to get generation, retrying in 1 second... (error: %v)\n", i+1, getErr)
		time.Sleep(1 * time.Second)
	}

	if getErr != nil {
		fmt.Printf("GetGeneration error after retries: %v\n", getErr)
		return
	}

	if genResp == nil {
		fmt.Println("Failed to retrieve generation data, response was nil.")
		return
	}

	fmt.Printf("Successfully retrieved generation:\n")
	fmt.Printf("  ID: %s\n", genResp.ID)
	fmt.Printf("  Model: %s\n", genResp.Model)
	fmt.Printf("  Prompt Tokens: %d\n", genResp.PromptTokens)
	fmt.Printf("  Completion Tokens: %d\n", genResp.CompletionTokens)
	fmt.Printf("  Total Tokens: %d\n", genResp.PromptTokens+genResp.CompletionTokens)
	fmt.Printf("  Cost: %f\n", genResp.TotalCost)
}
``` 
``` 