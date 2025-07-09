package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}
	folderID := os.Getenv("COZE_FOLDER_ID")

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	// Retrieve folders
	retrieveResp, err := client.Folders.Retrieve(ctx, &coze.RetrieveFolderReq{
		FolderID: folderID,
	})
	if err != nil {
		log.Fatalf("Retrieve folders failed: %v", err)
		return
	}

	fmt.Printf("Folder: %+v\n", retrieveResp)
}
