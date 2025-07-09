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
	workspaceID := os.Getenv("COZE_WORKSPACE_ID")
	folderType := os.Getenv("COZE_FOLDER_TYPE")

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	// List folders
	listResp, err := client.Folders.List(ctx, &coze.ListFoldersReq{
		WorkspaceID: workspaceID,
		FolderType:  coze.FolderType(folderType),
		PageSize:    10,
		PageNum:     1,
	})
	if err != nil {
		log.Fatalf("List folders failed: %v", err)
		return
	}

	for listResp.HasMore() {
		item := listResp.Current()
		fmt.Printf("Folder: %+v\n", item)
	}
	err = listResp.Err()
	if err != nil {
		log.Fatalf("List folders failed: %v", err)
		return
	}
}
