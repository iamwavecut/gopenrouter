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
