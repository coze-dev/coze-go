package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()

	spaceID := os.Getenv("WORKSPACE_ID")

	// Create list request with pagination parameters
	req := &coze.ListDatasetsReq{
		SpaceID:  spaceID,
		PageSize: 2,
		PageNum:  1,
	}

	// Get paginated results
	datasets, err := client.Datasets.List(ctx, req)
	if err != nil {
		fmt.Printf("Failed to list datasets: %v\n", err)
		return
	}

	for datasets.Next() {
		fmt.Println(datasets.Current())
	}
	if datasets.Err() != nil {
		fmt.Println("Error fetching datasets:", datasets.Err())
		return
	}

	// the page result will return followed information
	fmt.Println("total:", datasets.Total())
	fmt.Println("has_more:", datasets.HasMore())
}
