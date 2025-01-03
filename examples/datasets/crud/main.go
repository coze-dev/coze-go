package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth
	token := os.Getenv("COZE_API_TOKEN")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	spaceID := os.Getenv("WORKSPACE_ID")
	ctx := context.Background()

	/*
	 * Create dataset
	 */
	createReq := &coze.CreateDatasetsReq{
		Name:        "test dataset",
		SpaceID:     spaceID,
		FormatType:  coze.DocumentFormatTypeDocument,
		Description: "test dataset description",
	}

	createResp, err := client.Datasets.Create(ctx, createReq)
	if err != nil {
		fmt.Printf("Failed to create dataset: %v\n", err)
		return
	}
	fmt.Printf("Created dataset ID: %s\n", createResp.DatasetID)

	datasetID := createResp.DatasetID

	// Wait for 5 seconds
	time.Sleep(5 * time.Second)

	/*
	 * Update dataset
	 */
	updateReq := &coze.UpdateDatasetsReq{
		DatasetID:   datasetID,
		Name:        "updated dataset name",
		Description: "updated dataset description",
	}

	updateResp, err := client.Datasets.Update(ctx, updateReq)
	if err != nil {
		fmt.Printf("Failed to update dataset: %v\n", err)
		return
	}
	fmt.Printf("Update dataset response: %+v\n", updateResp)

	/*
	 * Delete dataset
	 */
	deleteReq := &coze.DeleteDatasetsReq{
		DatasetID: datasetID,
	}

	deleteResp, err := client.Datasets.Delete(ctx, deleteReq)
	if err != nil {
		fmt.Printf("Failed to delete dataset: %v\n", err)
		return
	}
	fmt.Printf("Delete dataset response: %+v\n", deleteResp)
}
