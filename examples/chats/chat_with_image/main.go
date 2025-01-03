package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/coze-dev/coze-go"
)

// This examples is about how to use the streaming interface to start a chats request
// with image upload and handle chats events
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	botID := os.Getenv("PUBLISHED_BOT_ID")
	userID := os.Getenv("USER_ID")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()

	// Call the upload file interface to get the image id.
	images, err := os.Open(os.Getenv("IMAGE_FILE_PATH"))
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	imageInfo, err := cozeCli.Files.Upload(ctx, &coze.UploadFilesReq{File: images})
	if err != nil {
		fmt.Println("Error uploading image:", err)
		return
	}
	fmt.Printf("upload image success, image id:%s\n", imageInfo.FileInfo.ID)

	//
	// Step one, create chats
	// Call the coze.Create.Stream() method to create a chats. The create method is a streaming
	// chats and will return a channel of ChatEvent. Developers should iterate the channel to get
	// chats events and handle them.

	req := &coze.CreateChatsReq{
		BotID:  botID,
		UserID: userID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionObjects([]*coze.MessageObjectString{
				coze.NewTextMessageObject("Describe this picture"),
				coze.NewImageMessageObjectByID(imageInfo.FileInfo.ID),
			}, nil),
		},
	}

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
		} else {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("done, log:%s\n", resp.Response().LogID())
}
