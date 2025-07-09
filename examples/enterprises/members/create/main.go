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

	fmt.Println("Creating enterprise members...")
	createReq := &coze.CreateEnterpriseMemberReq{
		EnterpriseID: enterpriseID,
		Users: []*coze.EnterpriseMember{
			{
				UserID: userID,
				Role:   coze.EnterpriseMemberRoleMember,
			},
		},
	}
	createResp, err := client.Enterprises.Members.Create(ctx, createReq)
	if err != nil {
		fmt.Printf("Failed to create enterprise members: %v\n", err)
		return
	}
	fmt.Printf("Created enterprise members successfully - Log ID: %s\n", createResp.LogID())
}
