package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/coze-dev/coze-go"
	"github.com/coze-dev/coze-go/examples/websockets/util"
	"github.com/gorilla/websocket"
)

type handler struct {
	coze.BaseWebSocketChatHandler
	data []byte
}

func (r *handler) OnClientError(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketClientErrorEvent) error {
	if errors.Is(event.Data, net.ErrClosed) {
		return nil
	}
	var wsErr *websocket.CloseError
	if errors.As(event.Data, &wsErr) {
		return nil
	}
	fmt.Printf("chat client_error=%s\n", event.Data)
	return nil
}

func (r *handler) OnError(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketErrorEvent) error {
	fmt.Printf("chat error=%s\n", event.Data)
	return nil
}

func (r *handler) OnConversationMessageDelta(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationMessageDeltaEvent) error {
	fmt.Printf("chat message_delta=%s\n", event.Data.Content)
	return nil
}

func (r *handler) OnConversationAudioDelta(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationAudioDeltaEvent) error {
	r.data = append(r.data, event.Data.Content...)
	return nil
}

func (r *handler) OnConversationAudioCompleted(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationAudioCompletedEvent) error {
	err := util.WritePCMToWavFile("output_chat.wav", r.data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("chat completed, audio write to %s", "output_chat.wav")
	return nil
}

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}
	botID := os.Getenv("COZE_BOT_ID")

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli,
		coze.WithBaseURL(cozeAPIBase),
		coze.WithLogLevel(coze.LogLevelDebug),
	)

	// Create chat WebSocket client
	chatClient := client.WebSockets.Chat.Create(context.Background(), &coze.CreateWebsocketChatReq{
		BotID: &botID,
	})
	chatClient.RegisterHandler(&handler{})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer chatClient.Close()

	// Send a message
	message := "今天天气真不错"
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
	if err := chatClient.Wait(); err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Chat completed!\n")
}
