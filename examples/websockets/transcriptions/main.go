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

	// Get base URL from environment (default to CN)
	baseURL := os.Getenv("COZE_API_BASE")
	if baseURL == "" {
		baseURL = coze.CnBaseURL
	}

	// Create Coze API client
	auth := coze.NewTokenAuth(apiToken)
	client := coze.NewCozeAPI(auth, coze.WithBaseURL(baseURL))

	// Create transcriptions WebSocket client
	transcriptionsClient := client.WebSockets.Audio.Transcriptions()

	// Set up event handlers
	handler := &coze.TranscriptionsEventHandler{
		OnTranscriptionsCreated: func(event coze.IWebSocketEvent) error {
			fmt.Println("Transcriptions session created")
			return nil
		},
		OnTranscriptionsMessageUpdate: func(event *coze.TranscriptionsMessageUpdateEvent) error {
			fmt.Printf("Transcription result: %s\n", event.Data.Content)
			return nil
		},
		OnTranscriptionsMessageCompleted: func(event coze.IWebSocketEvent) error {
			fmt.Println("Transcription completed")
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
	handler.RegisterHandlers(transcriptionsClient)

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := transcriptionsClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer transcriptionsClient.Close()

	// Wait for connection to be established
	time.Sleep(1 * time.Second)

	// Simulate sending audio data (in a real implementation, this would be actual audio data)
	// For this example, we'll just send some dummy data
	fmt.Println("Sending audio data...")
	audioData := []byte("This is simulated audio data for transcription testing")

	if err := transcriptionsClient.AppendAudioBuffer(audioData); err != nil {
		log.Fatalf("Failed to append audio: %v", err)
	}

	if err := transcriptionsClient.CompleteAudioBuffer(); err != nil {
		log.Fatalf("Failed to complete audio buffer: %v", err)
	}

	// Wait for transcription completion
	fmt.Println("Waiting for transcription completion...")
	event, err := transcriptionsClient.WaitForTranscriptionCompleted(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Transcription completed! Event: %+v\n", event)
}
