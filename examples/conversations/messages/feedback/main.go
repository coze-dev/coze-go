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

	conversationID := os.Getenv("CONVERSATION_ID")
	messageID := os.Getenv("MESSAGE_ID")
	ctx := context.Background()

	/*
	 * Add feedback to a message
	 */
	fmt.Println("Creating feedback...")
	createResp, err := client.Conversations.Messages.Feedback.Create(ctx, &coze.CreateConversationMessageFeedbackReq{
		ConversationID: conversationID,
		MessageID:      messageID,
		FeedbackType:   coze.FeedbackTypeLike,
		ReasonTypes:    []string{"helpful", "accurate"},
		Comment:        "This response was very helpful!",
	})
	if err != nil {
		fmt.Printf("Failed to create feedback: %v\n", err)
		return
	}
	fmt.Printf("Feedback created successfully - LogID: %s\n", createResp.LogID())

	/*
	 * Delete feedback from a message
	 */
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