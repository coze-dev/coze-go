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

	templateID := os.Getenv("COZE_TEMPLATE_ID")
	workspaceID := os.Getenv("WORKSPACE_ID")
	ctx := context.Background()

	/*
	 * Duplicate template
	 */
	newName := "duplicated_template"
	req := &coze.DuplicateTemplateReq{
		WorkspaceID: workspaceID,
		Name:        &newName,
	}

	resp, err := client.Templates.Duplicate(ctx, templateID, req)
	if err != nil {
		fmt.Printf("Failed to duplicate template: %v\n", err)
		return
	}
	fmt.Printf("Duplicated template - Entity ID: %s, Entity Type: %s\n", resp.EntityID, resp.EntityType)
}
