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

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	workspaceID := os.Getenv("WORKSPACE_ID")
	ctx := context.Background()

	/*
	 * List workflows
	 */
	workflowMode := coze.WorkflowModeWorkflow
	req := &coze.ListWorkflowsReq{
		WorkspaceID:  workspaceID,
		WorkflowMode: &workflowMode,
		PageNum:      1,
		PageSize:     10,
	}

	workflowsPaged, err := client.Workflows.List(ctx, req)
	if err != nil {
		fmt.Printf("Failed to list workflows: %v\n", err)
		return
	}

	fmt.Printf("Listing workflows:\n")
	for workflowsPaged.Next() {
		workflow := workflowsPaged.Current()
		fmt.Printf("- ID: %s, Name: %s, Description: %s\n", 
			workflow.WorkflowID, workflow.WorkflowName, workflow.Description)
	}

	if workflowsPaged.Err() != nil {
		fmt.Printf("Error during pagination: %v\n", workflowsPaged.Err())
	}

	/*
	 * List workflows with filters
	 */
	publishStatus := coze.PublishStatusPublished
	filteredReq := &coze.ListWorkflowsReq{
		WorkspaceID:   workspaceID,
		WorkflowMode:  &workflowMode,
		PublishStatus: &publishStatus,
		PageNum:       1,
		PageSize:      5,
	}

	filteredWorkflows, err := client.Workflows.List(ctx, filteredReq)
	if err != nil {
		fmt.Printf("Failed to list filtered workflows: %v\n", err)
		return
	}

	fmt.Printf("\nListing published workflows:\n")
	for filteredWorkflows.Next() {
		workflow := filteredWorkflows.Current()
		fmt.Printf("- ID: %s, Name: %s, App ID: %s\n", 
			workflow.WorkflowID, workflow.WorkflowName, workflow.AppID)
	}

	if filteredWorkflows.Err() != nil {
		fmt.Printf("Error during filtered pagination: %v\n", filteredWorkflows.Err())
	}
}