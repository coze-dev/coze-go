package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Bearer
	// Get an access_token through personal access token or oauth.
	token := "pat_83SdH03wfphMf2DpF2zsLg0mB8PxKAzTTADN6do3zXAr7SKlevyEuxAYn4BWOgSA"
	workflowID := "7485008793253789730"
	botID := "7483457315216572453"
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(coze.CnBaseURL))

	//
	// Step one, create chats
	req := &coze.WorkflowsChatStreamReq{
		BotID:      &botID,
		WorkflowID: workflowID,
		AdditionalMessages: []*coze.Message{
			coze.BuildUserQuestionText("你好", nil),
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
		if event.Event == coze.ChatEventDone {
			fmt.Println(event.WorkflowDebug.DebugUrl)
		} else if event.Event == coze.ChatEventConversationMessageDelta {
			fmt.Print(event.Message.Content)
		} else if event.Event == coze.ChatEventConversationChatCompleted {
			fmt.Printf("Token usage:%d\n", event.Chat.Usage.TokenCount)
		} else {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("done, log:%s\n", resp.Response().LogID())
}
