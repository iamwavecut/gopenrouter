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
