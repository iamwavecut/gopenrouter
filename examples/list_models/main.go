package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iamwavecut/gopenrouter"
	catalogapi "github.com/iamwavecut/gopenrouter/catalog"
)

func main() {
	client := gopenrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))
	api := catalogapi.New(client)
	ctx := context.Background()

	models, err := api.ListModels(ctx)
	if err != nil {
		fmt.Printf("catalog.ListModels error: %v\n", err)
		return
	}

	fmt.Println("Available Models:")
	for _, model := range models.Data {
		fmt.Printf("- ID: %s, Name: %s\n", model.ID, model.Name)
	}
}
