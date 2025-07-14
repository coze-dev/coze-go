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

func (r *handler) OnConversationAudioDelta(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationAudioDeltaEvent) error {
	r.data = append(r.data, event.Data.Content...)
	return nil
}

func (r *handler) OnConversationAudioCompleted(ctx context.Context, cli *coze.WebSocketChat, event *coze.WebSocketConversationAudioCompletedEvent) error {
	filename := "output_input_text_generate_audio.wav"
	err := util.WritePCMToWavFile(filename, r.data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("generate audio completed, audio write to %s", filename)
	return nil
}

// This example demonstrates how to use the input_text.generate_audio event
// for voice synthesis in WebSocket chat
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
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer chatClient.Close()

	// Example: Use input_text.generate_audio to synthesize voice without triggering bot response
	fmt.Println("Generating voice synthesis without bot response...")
	if err := chatClient.InputTextGenerateAudio(&coze.WebSocketInputTextGenerateAudioEventData{
		Mode: coze.WebSocketInputTextGenerateAudioModeText,
		Text: "亲，你怎么不说话了。",
	}); err != nil {
		log.Fatalf("Failed to generate audio: %v", err)
	}

	// Wait for the chat to complete
	if err := chatClient.Wait(
		coze.WebSocketEventTypeConversationAudioCompleted,
		coze.WebSocketEventTypeError,
	); err != nil {
		log.Fatalf("Chat completed with error: %v", err)
	}
}
