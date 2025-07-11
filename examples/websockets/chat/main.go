package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get API token from environment
	apiToken := os.Getenv("COZE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("COZE_API_TOKEN environment variable is required")
	}

	// Get bot ID from environment
	botID := os.Getenv("COZE_BOT_ID")
	if botID == "" {
		log.Fatal("COZE_BOT_ID environment variable is required")
	}

	// Get base URL from environment (default to CN)
	baseURL := os.Getenv("COZE_API_BASE")
	if baseURL == "" {
		baseURL = coze.CnBaseURL
	}

	// Create Coze API client
	auth := coze.NewTokenAuth(apiToken)
	client := coze.NewCozeAPI(auth, coze.WithBaseURL(baseURL))

	// Create chat WebSocket client
	chatClient := client.WebSockets.Chat.Create(context.Background(), &coze.CreateWebsocketChatReq{
		BotID: &botID,
	})

	chatClient.OnChatCreated(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketChatCreatedEvent) error {
		fmt.Println("Chat session created")
		return nil
	})
	chatClient.OnConversationChatCreated(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationChatCreatedEvent) error {
		fmt.Printf("Conversation chat created: %s\n", event.Data.ID)
		return nil
	})
	chatClient.OnConversationMessageDelta(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationMessageDeltaEvent) error {
		fmt.Printf("Message delta: %s\n", event.Data.Content)
		return nil
	})
	chatClient.OnConversationChatCompleted(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationChatCompletedEvent) error {
		fmt.Printf("Chat completed: %s\n", event.Data.ID)
		return nil
	})
	chatClient.OnError(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketErrorEvent) error {
		fmt.Printf("Error: %v\n", event)
		return nil
	})
	chatClient.OnClosed(func(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketClosedEvent) error {
		fmt.Println("Connection closed")
		return nil
	})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer chatClient.Close()

	// Send a message
	message := "Hello! How are you today?"
	fmt.Printf("Sending message: %s\n", message)
	if err := chatClient.ConversationMessageCreate(&coze.WebSocketConversationMessageCreateEventData{
		Role:        coze.MessageRoleUser,
		ContentType: coze.MessageContentTypeText,
		Content:     message,
	}); err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}

	// Wait for chat completion
	fmt.Println("Waiting for chat completion...")
	err := chatClient.Wait()
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Chat completed!\n")
}
