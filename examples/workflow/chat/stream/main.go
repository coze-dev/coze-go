package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func ofNotZero[T comparable](v T) *T {
	var empty T
	if v == empty {
		return &v
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	workflowID := os.Getenv("COZE_WORKFLOW_ID")
	botID := os.Getenv("COZE_BOT_ID")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	//
	// Step one, create chats
	req := &coze.WorkflowsChatStreamReq{
		BotID:      ofNotZero(botID),
		WorkflowID: workflowID,
		AdditionalMessages: []*coze.Message{
			coze.BuildUserQuestionText("What can you do?", nil),
		},
		Parameters: map[string]any{
			"name": "John",
		},
	}

	resp, err := cozeCli.Workflows.Chat.Stream(ctx, req)
	if err != nil {
		fmt.Printf("Error starting chats: %v\n", err)
		return
	}

	defer resp.Close()
	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if event.Event == coze.ChatEventConversationMessageDelta {
			fmt.Print(event.Message.Content)
		} else if event.Event == coze.ChatEventConversationChatCompleted {
			fmt.Printf("\n")
			fmt.Printf("Token usage:%d\n", event.Chat.Usage.TokenCount)
		} else if event.Event == coze.ChatEventDone {
			fmt.Printf("debug_url: %s\n", event.WorkflowDebug.DebugUrl)
		}
	}

	fmt.Printf("done, log:%s\n", resp.Response().LogID())
}
