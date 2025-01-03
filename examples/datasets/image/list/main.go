package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth
	token := os.Getenv("COZE_API_TOKEN")

	// Initialize the Coze client with custom configuration
	client := coze.NewCozeAPI(
		coze.NewTokenAuth(token),
		coze.WithBaseURL(os.Getenv("COZE_API_BASE")),
	)
	ctx := context.Background()

	datasetID := os.Getenv("DATASET_ID")

	// Create list request with pagination parameters
	req := &coze.ListDatasetsImagesReq{
		DatasetID: datasetID,
		PageSize:  2,
		PageNum:   1,
	}

	// Get paginated results
	images, err := client.Datasets.Images.List(ctx, req)
	if err != nil {
		fmt.Printf("Failed to list images: %v\n", err)
		return
	}

	for images.Next() {
		fmt.Println(images.Current())
	}
	if images.Err() != nil {
		fmt.Println("Error fetching images:", images.Err())
		return
	}

	// the page result will return followed information
	fmt.Println("total:", images.Total())
	fmt.Println("has_more:", images.HasMore())
}
