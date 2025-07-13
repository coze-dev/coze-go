package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coze-dev/coze-go"
	"github.com/coze-dev/coze-go/examples/websockets/util"
)

type handler struct {
	coze.BaseWebSocketAudioSpeechHandler
	data []byte
}

func (r *handler) OnClientError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketClientErrorEvent) error {
	log.Printf("speech client_error: %v", event)
	return nil
}

func (r *handler) OnError(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketErrorEvent) error {
	log.Printf("speech error: %v", event)
	return nil
}

func (r *handler) OnSpeechAudioUpdate(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioUpdateEvent) error {
	r.data = append(r.data, event.Data.Delta...)
	return nil
}

func (r *handler) OnSpeechAudioCompleted(ctx context.Context, cli *coze.WebSocketAudioSpeech, event *coze.WebSocketSpeechAudioCompletedEvent) error {
	filename := "output_speech.wav"
	err := util.WritePCMToWavFile(filename, r.data)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return err
	}

	fmt.Printf("speech completed, audio write to %s\n", filename)
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
	client := coze.NewCozeAPI(authCli,
		coze.WithBaseURL(cozeAPIBase),
		coze.WithLogLevel(coze.LogLevelDebug),
	)

	// Create speech WebSocket client
	speechClient := client.WebSockets.Audio.Speech.Create(context.Background(), &coze.CreateWebsocketAudioSpeechReq{})
	// speechClient.RegisterHandler(&handler{})

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := speechClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer speechClient.Close()

	// Send text to be converted to speech
	text := "今天天气不错"
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
	if err := speechClient.Wait(); err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}
}
