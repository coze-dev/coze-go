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

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	// List apps
	listResp, err := client.Apps.List(ctx, &coze.ListAppReq{
		WorkspaceID: workspaceID,
		PageSize:    10,
		PageNum:     1,
	})
	if err != nil {
		log.Fatalf("List apps failed: %v", err)
		return
	}

	for listResp.HasMore() {
		item := listResp.Current()
		fmt.Printf("App: %+v\n", item)
	}
	err = listResp.Err()
	if err != nil {
		log.Fatalf("List apps failed: %v", err)
		return
	}
}
