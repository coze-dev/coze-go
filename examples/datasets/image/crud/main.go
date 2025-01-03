package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth
	token := os.Getenv("COZE_API_TOKEN")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	datasetID := os.Getenv("DATASET_ID")
	ctx := context.Background()
	/*
	 * Create image document
	 */

	/*
	 * Step 1: Upload image to Coze
	 */
	imagePath := os.Getenv("IMAGE_FILE_PATH")
	fileInfo, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("Failed to open image file: %v\n", err)
		return
	}
	imageInfo, err := client.Files.Upload(ctx, &coze.UploadFilesReq{
		File: fileInfo,
	})
	if err != nil {
		fmt.Printf("Failed to upload image: %v\n", err)
		return
	}
	fmt.Printf("Image uploaded: %+v\n", imageInfo)

	/*
	 * Step 2: Create document
	 */
	datasetIDInt, _ := strconv.ParseInt(datasetID, 10, 64)
	fileIDInt, _ := strconv.ParseInt(imageInfo.FileInfo.ID, 10, 64)

	createReq := &coze.CreateDatasetsDocumentsReq{
		DatasetID: datasetIDInt,
		DocumentBases: []*coze.DocumentBase{
			coze.DocumentBaseBuildImage("test image", fileIDInt),
		},
		FormatType: coze.DocumentFormatTypeImage,
	}

	createResp, err := client.Datasets.Documents.Create(ctx, createReq)
	if err != nil {
		fmt.Printf("Failed to create document: %v\n", err)
		return
	}
	fmt.Printf("Document created: %+v\n", createResp)

	/*
	 * Step 3: Make sure upload is completed
	 */
	documentID := createResp.DocumentInfos[0].DocumentID

	processReq := &coze.ProcessDocumentsReq{
		DatasetID:   datasetID,
		DocumentIDs: []string{documentID},
	}

	// Poll until processing is complete
	for {
		processResp, err := client.Datasets.Process(ctx, processReq)
		if err != nil {
			fmt.Printf("Failed to check process status: %v\n", err)
			return
		}
		fmt.Printf("Process status: %+v\n", processResp)

		status := processResp.Data[0].Status
		if status == coze.DocumentStatusProcessing {
			fmt.Printf("Upload is not completed, please wait, process: %d\n", processResp.Data[0].Progress)
			time.Sleep(time.Second)
		} else if status == coze.DocumentStatusFailed {
			fmt.Println("Upload failed, please check")
			return
		} else {
			fmt.Println("Upload completed")
			break
		}
	}

	/*
	 * Update image caption
	 */
	caption := "new image caption"
	updateReq := &coze.UpdateDatasetImageReq{
		DatasetID:  datasetID,
		DocumentID: documentID,
		Caption:    &caption,
	}

	updateResp, err := client.Datasets.Images.Update(ctx, updateReq)
	if err != nil {
		fmt.Printf("Failed to update image: %v\n", err)
		return
	}
	fmt.Printf("Image updated: %+v\n", updateResp)
}
