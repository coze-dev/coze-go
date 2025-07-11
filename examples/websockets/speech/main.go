package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

type handler struct {
}

func (r handler) OnClientError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketClientErrorEvent) error {
	log.Printf("speech client_error: %v", event)
	return nil
}

func (r handler) OnClosed(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketClosedEvent) error {
	log.Printf("speech closed")
	return nil
}

func (r handler) OnError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketErrorEvent) error {
	log.Printf("speech error: %v", event)
	return nil
}

func (r handler) OnSpeechCreated(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechCreatedEvent) error {
	log.Printf("speech created: %v", event)
	return nil
}

func (r handler) OnSpeechUpdated(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechUpdatedEvent) error {
	log.Printf("speech updated: %v", event)
	return nil
}

func (r handler) OnInputTextBufferCompleted(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketInputTextBufferCompletedEvent) error {
	log.Printf("speech input_text_buffer_completed: %v", event)
	return nil
}

func (r handler) OnSpeechAudioUpdate(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioUpdateEvent) error {
	log.Printf("speech audio_update: %v", event)
	return nil
}

func (r handler) OnSpeechAudioCompleted(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioCompletedEvent) error {
	log.Printf("speech audio_completed: %v", event)
	return nil
}

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))

	// Create speech WebSocket client
	speechClient := client.WebSockets.Audio.Speech.Create(context.Background(), &coze.CreateWebsocketAudioSpeechReq{
		// Model: "tts-1",
	})
	speechClient.RegisterHandler(&handler{})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := speechClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer speechClient.Close()

	// Send text to be converted to speech
	text := "Hello, this is a test of the WebSocket speech functionality!"
	fmt.Printf("Sending text: %s\n", text)

	if err := speechClient.InputTextBufferAppend(&coze.WebSocketInputTextBufferAppendEventData{
		Delta: text,
	}); err != nil {
		log.Fatalf("Failed to append text: %v", err)
	}

	if err := speechClient.InputTextBufferComplete(nil); err != nil {
		log.Fatalf("Failed to complete text buffer: %v", err)
	}

	// Wait for speech completion
	fmt.Println("Waiting for speech completion...")
	event, err := speechClient.Wait(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Speech completed! Event: %+v\n", event)
}
