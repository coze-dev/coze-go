package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
	"github.com/coze-dev/coze-go/websockets"
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

	// Create speech WebSocket client
	speechClient := client.WebSockets.Audio.Speech()

	// Set up event handlers
	handler := &websockets.SpeechEventHandler{
		OnSpeechCreated: func(event *websockets.SpeechCreatedEvent) error {
			fmt.Printf("Speech session created: %s\n", event.Data.SessionID)
			return nil
		},
		OnSpeechAudioUpdate: func(event *websockets.SpeechAudioUpdateEvent) error {
			fmt.Printf("Received audio data chunk (length: %d)\n", len(event.Data.Delta))
			return nil
		},
		OnSpeechAudioCompleted: func(event *websockets.SpeechAudioCompletedEvent) error {
			fmt.Printf("Speech audio completed: %s\n", event.Data.SessionID)
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
	handler.RegisterHandlers(speechClient)

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := speechClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer speechClient.Close()

	// Wait for connection to be established
	time.Sleep(1 * time.Second)

	// Send text to be converted to speech
	text := "Hello, this is a test of the WebSocket speech functionality!"
	fmt.Printf("Sending text: %s\n", text)
	
	if err := speechClient.AppendTextBuffer(text); err != nil {
		log.Fatalf("Failed to append text: %v", err)
	}

	if err := speechClient.CompleteTextBuffer(); err != nil {
		log.Fatalf("Failed to complete text buffer: %v", err)
	}

	// Wait for speech completion
	fmt.Println("Waiting for speech completion...")
	event, err := speechClient.WaitForSpeechAudioCompleted(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Speech completed! Event: %+v\n", event)
}