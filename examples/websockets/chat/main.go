package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	chatClient := client.WebSockets.Chat.Create(
		coze.WithBotID(botID),
	)

	// Set up event handlers
	handler := &coze.ChatEventHandler{
		OnChatCreated: func(event *coze.WebSocketEvent) error {
			fmt.Println("Chat session created")
			return nil
		},
		OnConversationChatCreated: func(event *coze.ConversationChatCreatedEvent) error {
			fmt.Printf("Conversation chat created: %s\n", event.Data.ChatID)
			return nil
		},
		OnConversationMessageDelta: func(event *coze.ConversationMessageDeltaEvent) error {
			fmt.Printf("Message delta: %s\n", event.Data.Content)
			return nil
		},
		OnConversationAudioDelta: func(event *coze.ConversationAudioDeltaEvent) error {
			audioData := event.Data.GetAudio()
			fmt.Printf("Audio delta received (length: %d)\n", len(audioData))
			return nil
		},
		OnConversationChatCompleted: func(event *coze.ConversationChatCompletedEvent) error {
			fmt.Printf("Chat completed: %s\n", event.Data.ChatID)
			return nil
		},
		OnConversationChatRequiresAction: func(event *coze.ConversationChatRequiresActionEvent) error {
			fmt.Printf("Chat requires action: %s\n", event.Data.ChatID)

			// Example: Handle tool call requirements
			if event.Data.RequiredAction != nil && event.Data.RequiredAction.SubmitToolOutputs != nil {
				for _, toolCall := range event.Data.RequiredAction.SubmitToolOutputs.ToolCalls {
					fmt.Printf("Tool call required: %s (%s)\n", toolCall.Function.Name, toolCall.ID)

					// In a real implementation, you would execute the tool and get the result
					// For this example, we'll just return a dummy result
					toolOutputs := []coze.ToolOutput{
						{
							ToolCallID: toolCall.ID,
							Output:     "Tool execution result: success",
						},
					}

					// Submit tool outputs
					if err := chatClient.SubmitToolOutputs(event.Data.ChatID, toolOutputs); err != nil {
						fmt.Printf("Failed to submit tool outputs: %v\n", err)
					}
				}
			}
			return nil
		},
		OnError: func(err error) error {
			fmt.Printf("Error: %v\n", err)
			return nil
		},
		OnClosed: func() error {
			fmt.Println("Connection closed")
			return nil
		},
	}

	// Register event handlers
	handler.RegisterHandlers(chatClient)

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer chatClient.Close()

	// Wait for connection to be established
	time.Sleep(1 * time.Second)

	// Update chat configuration
	fmt.Println("Updating chat configuration...")
	if err := chatClient.UpdateChat(botID, nil, nil); err != nil {
		log.Fatalf("Failed to update chat: %v", err)
	}

	// Send a message
	message := "Hello! How are you today?"
	fmt.Printf("Sending message: %s\n", message)

	if err := chatClient.CreateMessage(message); err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}

	// Alternative: Send audio data
	// audioData := []byte("This is simulated audio data")
	// if err := chatClient.AppendAudioBuffer(audioData); err != nil {
	//     log.Fatalf("Failed to append audio: %v", err)
	// }
	// if err := chatClient.CompleteAudioBuffer(); err != nil {
	//     log.Fatalf("Failed to complete audio buffer: %v", err)
	// }

	// Wait for chat completion
	fmt.Println("Waiting for chat completion...")
	event, err := chatClient.WaitForChatCompleted(60 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Chat completed! Event: %+v\n", event)
}
