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
