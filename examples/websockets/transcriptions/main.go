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

func (r handler) OnClientError(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketClientErrorEvent) error {
	fmt.Printf("transcriptions client error: %v", event)
	return nil
}

func (r handler) OnClosed(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketClosedEvent) error {
	fmt.Printf("transcriptions closed: %v", event)
	return nil
}

func (r handler) OnError(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketErrorEvent) error {
	fmt.Printf("transcriptions error: %v", event)
	return nil
}

func (r handler) OnTranscriptionsCreated(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketTranscriptionsCreatedEvent) error {
	fmt.Printf("transcriptions created: %v", event)
	return nil
}

func (r handler) OnTranscriptionsUpdated(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketTranscriptionsUpdatedEvent) error {
	fmt.Printf("transcriptions updated: %v", event)
	return nil
}

func (r handler) OnInputAudioBufferCompleted(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketInputAudioBufferCompletedEvent) error {
	fmt.Printf("transcriptions input audio buffer completed: %v", event)
	return nil
}

func (r handler) OnInputAudioBufferCleared(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketInputAudioBufferClearedEvent) error {
	fmt.Printf("transcriptions input audio buffer cleared: %v", event)
	return nil
}

func (r handler) OnTranscriptionsMessageUpdate(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketTranscriptionsMessageUpdateEvent) error {
	fmt.Printf("transcriptions message update: %v", event)
	return nil
}

func (r handler) OnTranscriptionsMessageCompleted(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketTranscriptionsMessageCompletedEvent) error {
	fmt.Printf("transcriptions message completed: %v", event)
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

	// Create transcriptions WebSocket client
	transcriptionsClient := client.WebSockets.Audio.Transcriptions.Create(context.Background(), &coze.CreateWebsocketAudioTranscriptionReq{})
	transcriptionsClient.RegisterHandler(&handler{})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := transcriptionsClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer transcriptionsClient.Close()

	// Simulate sending audio data (in a real implementation, this would be actual audio data)
	// For this example, we'll just send some dummy data
	fmt.Println("Sending audio data...")
	audioData := []byte("This is simulated audio data for transcription testing")

	if err := transcriptionsClient.InputAudioBufferAppend(&coze.WebSocketInputAudioBufferAppendEventData{
		Delta: audioData,
	}); err != nil {
		log.Fatalf("Failed to append audio: %v", err)
	}

	if err := transcriptionsClient.InputAudioBufferComplete(nil); err != nil {
		log.Fatalf("Failed to complete audio buffer: %v", err)
	}

	// Wait for transcription completion
	fmt.Println("Waiting for transcription completion...")
	event, err := transcriptionsClient.Wait(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Transcription completed! Event: %+v\n", event)
}
