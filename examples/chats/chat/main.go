package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

// This examples describes how to use the chats interface to initiate conversations,
// poll the status of the conversation, and obtain the messages after the conversation is completed.

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	botID := os.Getenv("PUBLISHED_BOT_ID")
	uid := os.Getenv("USER_ID")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()

	//
	// Step one, create chats
	// Call the coze.Create.Create() method to create a chats. The create method is a non-streaming
	// chats and will return a Create class. Developers should periodically check the status of the
	// chats and handle them separately according to different states.
	//
	req := &coze.CreateChatsReq{
		BotID:  botID,
		UserID: uid,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText("What can you do?", nil),
		},
	}

	chatResp, err := cozeCli.Chat.Create(ctx, req)
	if err != nil {
		fmt.Println("Error creating chats:", err)
		return
	}
	fmt.Println(chatResp)
	fmt.Println(chatResp.LogID())
	chat := chatResp.Chat
	chatID := chat.ID
	conversationID := chat.ConversationID

	//
	// Step two, poll the result of chats
	// Assume the development allows at most one chats to runs for 10 seconds. If it exceeds 10 seconds,
	// the chats will be cancelled.
	// And when the chats status is not completed, poll the status of the chats once every second.
	// After the chats is completed, retrieve all messages in the chats.
	//
	timeout := time.After(1) // time.Second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for chat.Status == coze.ChatStatusInProgress {
		select {
		case <-timeout:
			// The chats can be cancelled before its completed.
			cancelResp, err := cozeCli.Chat.Cancel(ctx, &coze.CancelChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			})
			if err != nil {
				fmt.Println("Error cancelling chats:", err)
			}
			fmt.Println("cancel")
			fmt.Println(cancelResp)
			fmt.Println(cancelResp.LogID())
			break
		case <-ticker.C:
			resp, err := cozeCli.Chat.Retrieve(ctx, &coze.RetrieveChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			})
			if err != nil {
				fmt.Println("Error retrieving chats:", err)
				continue
			}
			fmt.Println("retrieve")
			fmt.Println(resp)
			fmt.Println(resp.LogID())
			chat = resp.Chat
			if chat.Status == coze.ChatStatusCompleted {
				break
			}
		}
	}

	// The sdk provide an automatic polling method.
	chat2, err := cozeCli.Chat.CreateAndPoll(ctx, req, nil)
	if err != nil {
		fmt.Println("Error in CreateAndPoll:", err)
		return
	}
	fmt.Println(chat2)

	// the developer can also set the timeout.
	pollTimeout := 10
	chat3, err := cozeCli.Chat.CreateAndPoll(ctx, req, &pollTimeout)
	if err != nil {
		fmt.Println("Error in CreateAndPollWithTimeout:", err)
		return
	}
	fmt.Println(chat3)
}
