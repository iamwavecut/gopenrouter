package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	logprobs := true
	top := 5
	req := gopenrouter.ChatCompletionRequest{
		Model: "openai/gpt-4o",
		Messages: []gopenrouter.ChatCompletionMessage{
			{Role: gopenrouter.RoleUser, Content: "Write a short haiku about wind."},
		},
		LogProbs:    &logprobs,
		TopLogProbs: &top,
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Logprobs == nil {
		fmt.Println(resp.Choices[0].Message.Content)
		return
	}

	fmt.Println("Tokens with logprobs:")
	for _, t := range resp.Choices[0].Logprobs.Content {
		fmt.Printf("%s\t%.3f\n", t.Token, t.LogProb)
		for _, alt := range t.TopLogProbs {
			fmt.Printf("  alt: %s\t%.3f\n", alt.Token, alt.LogProb)
		}
	}
}
