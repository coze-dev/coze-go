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

	// List workflows
	listResp, err := client.Workflows.List(ctx, &coze.ListWorkflowReq{
		WorkspaceID: &workspaceID,
		PageSize:    10,
		PageNum:     1,
	})
	if err != nil {
		log.Fatalf("List workflows failed: %v", err)
		return
	}

	for listResp.HasMore() {
		item := listResp.Current()
		fmt.Printf("Workflow: %+v\n", item)
	}
	err = listResp.Err()
	if err != nil {
		log.Fatalf("List workflows failed: %v", err)
		return
	}
}
