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
	enterpriseID := os.Getenv("COZE_ENTERPRISE_ID")
	userID := os.Getenv("COZE_USER_ID")
	receiverUserID := os.Getenv("COZE_RECEIVER_USER_ID")

	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

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
