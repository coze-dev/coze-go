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

	enterpriseID := os.Getenv("ENTERPRISE_ID")
	userID := os.Getenv("USER_ID")
	receiverUserID := os.Getenv("RECEIVER_USER_ID")
	ctx := context.Background()

	/*
	 * Create enterprise members
	 */
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

	/*
	 * Update enterprise member role
	 */
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

	/*
	 * Delete enterprise member
	 */
	fmt.Println("Deleting enterprise member...")
	deleteReq := &coze.DeleteEnterpriseMemberReq{
		EnterpriseID:   enterpriseID,
		UserID:         userID,
		ReceiverUserID: receiverUserID,
	}

	deleteResp, err := client.Enterprises.Members.Delete(ctx, deleteReq)
	if err != nil {
		fmt.Printf("Failed to delete enterprise member: %v\n", err)
		return
	}
	fmt.Printf("Deleted enterprise member successfully - Log ID: %s\n", deleteResp.LogID())
}
