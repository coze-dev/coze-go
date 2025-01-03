package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// This examples demonstrates how to get response log ID from different API calls.
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE"))) // coze.WithLogger() The developer can set the logger to print the log.
	// coze.WithLogLevel(log.LogDebug) The developer can set the log level.

	ctx := context.Background()
	botID := os.Getenv("COZE_BOT_ID")

	// Example 1: Get log ID from bot retrieve API
	botsResp, err := cozeCli.Bots.Retrieve(ctx, &coze.RetrieveBotsReq{
		BotID: botID,
	})
	if err != nil {
		fmt.Printf("Error retrieving bot: %v\n", err)
		return
	}
	fmt.Printf("Bot retrieve log ID: %s\n", botsResp.Response().LogID())

	// Example 2: Get log ID from chats API
	chatResp, err := cozeCli.Chat.Create(ctx, &coze.CreateChatsReq{
		BotID:  botID,
		UserID: os.Getenv("USER_ID"),
		Messages: []*coze.Message{
			coze.BuildUserQuestionText("What can you do?", nil),
		},
	})
	if err != nil {
		fmt.Printf("Error creating chats: %v\n", err)
		return
	}
	fmt.Printf("Create create log ID: %s\n", chatResp.Response().LogID())

	// Example 3: Get log ID from file upload API
	file, err := os.Open(os.Getenv("FILE_PATH"))
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	fileResp, err := cozeCli.Files.Upload(ctx, &coze.UploadFilesReq{File: file})
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		return
	}
	fmt.Printf("File upload log ID: %s\n", fileResp.Response().LogID())
}
