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

	conversationID := os.Getenv("COZE_CONVERSATION_ID")
	messageID := os.Getenv("COZE_MESSAGE_ID")

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	// Delete feedback from a message
	fmt.Println("Deleting feedback...")
	deleteResp, err := client.Conversations.Messages.Feedback.Delete(ctx, &coze.DeleteConversationMessageFeedbackReq{
		ConversationID: conversationID,
		MessageID:      messageID,
	})
	if err != nil {
		fmt.Printf("Failed to delete feedback: %v\n", err)
		return
	}
	fmt.Printf("Feedback deleted successfully - LogID: %s\n", deleteResp.LogID())
}
