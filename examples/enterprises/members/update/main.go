package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}

	enterpriseID := os.Getenv("ENTERPRISE_ID")
	userID := os.Getenv("USER_ID")

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	fmt.Println("Updating enterprise member role...")
	updateReq := &coze.UpdateEnterpriseMemberReq{
		EnterpriseID: enterpriseID,
		UserID:       userID,
		Role:         coze.EnterpriseMemberRoleAdmin,
	}

	updateResp, err := client.Enterprises.Members.Update(ctx, updateReq)
	if err != nil {
		fmt.Printf("Failed to update enterprise member: %v\n", err)
		return
	}
	fmt.Printf("Updated enterprise member successfully - Log ID: %s\n", updateResp.LogID())
}
