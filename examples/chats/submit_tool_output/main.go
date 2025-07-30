package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

// This use case teaches you how to use local plugin.
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	botID := os.Getenv("PUBLISHED_BOT_ID")
	userID := os.Getenv("USER_ID")

	authCli := coze.NewTokenAuth(token)

	customClient := &http.Client{
		Timeout: time.Minute * 20,
	}

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli,
		coze.WithBaseURL(os.Getenv("COZE_API_BASE")),
		coze.WithHttpClient(customClient),
	)

	ctx := context.Background()

	req := &coze.CreateChatsReq{
		BotID:  botID,
		UserID: userID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText("What's the weather like in Shenzhen today?", nil),
		},
	}

	var pluginEvent *coze.ChatEvent
	var conversationID string

	resp, err := cozeCli.Chat.Stream(ctx, req)
	if err != nil {
		fmt.Println("Error starting stream:", err)
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
			fmt.Printf("Token usage:%d\n", event.Chat.Usage.TokenCount)
		} else if event.Event == coze.ChatEventConversationChatRequiresAction {
			fmt.Println("need action")
			pluginEvent = event
			conversationID = event.Chat.ConversationID
			break
		} else {
			fmt.Printf("\n")
		}
	}

	if pluginEvent == nil {
		return
	}

	var toolOutputs []*coze.ToolOutput
	for _, callInfo := range pluginEvent.Chat.RequiredAction.SubmitToolOutputs.ToolCalls {
		callID := callInfo.ID
		// you can handle different plugin by name.
		functionName := callInfo.Function.Name
		// you should unmarshal arguments if necessary.
		argsJSON := callInfo.Function.Arguments

		fmt.Printf("Function called: %s with args: %s\n", functionName, argsJSON)
		toolOutputs = append(toolOutputs, &coze.ToolOutput{
			ToolCallID: callID,
			Output:     "It is 18 to 21",
		})
	}

	toolReq := &coze.SubmitToolOutputsChatReq{
		ChatID:         pluginEvent.Chat.ID,
		ConversationID: conversationID,
		ToolOutputs:    toolOutputs,
	}

	resp2, err := cozeCli.Chat.StreamSubmitToolOutputs(ctx, toolReq)
	if err != nil {
		fmt.Println("Error submitting tool outputs:", err)
		return
	}

	defer resp2.Close()
	for {
		event, err := resp2.Recv()
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
		} else {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("done, log:%s\n", resp.Response().LogID())
}
