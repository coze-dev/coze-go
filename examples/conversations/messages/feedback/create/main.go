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

	// Add feedback to a message
	fmt.Println("Creating feedback...")
	createResp, err := client.Conversations.Messages.Feedback.Create(ctx, &coze.CreateConversationMessageFeedbackReq{
		ConversationID: conversationID,
		MessageID:      messageID,
		FeedbackType:   coze.FeedbackTypeLike,
		ReasonTypes:    []string{"helpful", "accurate"},
		Comment:        &[]string{"This response was very helpful!"}[0],
	})
	if err != nil {
		fmt.Printf("Failed to create feedback: %v\n", err)
		return
	}
	fmt.Printf("Feedback created successfully - LogID: %s\n", createResp.LogID())
}
