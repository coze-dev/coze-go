package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coze-dev/coze-go"
)

type handler struct {
	coze.BaseWebSocketAudioTranscriptionHandler
}

func (r *handler) OnClientError(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketClientErrorEvent) error {
	fmt.Printf("transcriptions client error: %v\n", event)
	return nil
}

func (r *handler) OnError(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketErrorEvent) error {
	fmt.Printf("transcriptions error: %v\n", event)
	return nil
}

func (r *handler) OnTranscriptionsMessageUpdate(ctx context.Context, cli *coze.WebSocketAudioTranscription, event *coze.WebSocketTranscriptionsMessageUpdateEvent) error {
	fmt.Printf("transcriptions message update: %s\n", event.Data.Content)
	return nil
}

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}
	cozeAudioFile := os.Getenv("COZE_AUDIO_FILE")
	if cozeAudioFile == "" {
		cozeAudioFile = "output.wav"
	}

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli,
		coze.WithBaseURL(cozeAPIBase),
		coze.WithLogLevel(coze.LogLevelDebug),
	)

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
	audioData, err := os.ReadFile(cozeAudioFile)
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

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
	if err := transcriptionsClient.Wait(); err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Transcription completed!\n")
}
