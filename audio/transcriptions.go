package audio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-go/websockets"
)

// TranscriptionsClient handles audio transcriptions WebSocket connections
type TranscriptionsClient struct {
	wsClient *websockets.WebSocketClient
}

// TranscriptionsClientOption configures the transcriptions client
type TranscriptionsClientOption func(*TranscriptionsClient)

// WithInputAudio configures input audio settings
func WithInputAudio(inputAudio *websockets.InputAudio) TranscriptionsClientOption {
	return func(c *TranscriptionsClient) {
		// This will be sent during connection setup
	}
}

// NewTranscriptionsClient creates a new transcriptions WebSocket client
func NewTranscriptionsClient(baseURL string, auth websockets.Auth, opts ...TranscriptionsClientOption) *TranscriptionsClient {
	wsClient := websockets.NewWebSocketClient(
		baseURL,
		"v1/audio/transcriptions",
		auth,
	)
	
	client := &TranscriptionsClient{
		wsClient: wsClient,
	}
	
	for _, opt := range opts {
		opt(client)
	}
	
	return client
}

// Connect establishes the WebSocket connection
func (c *TranscriptionsClient) Connect() error {
	return c.wsClient.Connect()
}

// Close closes the WebSocket connection
func (c *TranscriptionsClient) Close() error {
	return c.wsClient.Close()
}

// IsConnected returns whether the client is connected
func (c *TranscriptionsClient) IsConnected() bool {
	return c.wsClient.IsConnected()
}

// UpdateTranscriptions updates the transcriptions configuration
func (c *TranscriptionsClient) UpdateTranscriptions(inputAudio *websockets.InputAudio) error {
	event := websockets.TranscriptionsUpdateEvent{
		EventType: websockets.EventTypeTranscriptionsUpdate,
		Data: &websockets.TranscriptionsUpdateData{
			InputAudio: inputAudio,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// AppendAudioBuffer appends audio data to the input buffer
func (c *TranscriptionsClient) AppendAudioBuffer(audioData []byte) error {
	// Encode audio data to base64
	encoded := base64.StdEncoding.EncodeToString(audioData)
	
	event := websockets.InputAudioBufferAppendEvent{
		EventType: websockets.EventTypeInputAudioBufferAppend,
		Data: &websockets.InputAudioBufferAppendData{
			Delta: encoded,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *TranscriptionsClient) CompleteAudioBuffer() error {
	event := websockets.InputAudioBufferCompleteEvent{
		EventType: websockets.EventTypeInputAudioBufferComplete,
	}
	
	return c.wsClient.SendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *TranscriptionsClient) ClearAudioBuffer() error {
	event := websockets.InputAudioBufferClearEvent{
		EventType: websockets.EventTypeInputAudioBufferClear,
	}
	
	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *TranscriptionsClient) OnEvent(eventType websockets.WebSocketEventType, handler websockets.EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForTranscriptionCompleted waits for transcription to complete
func (c *TranscriptionsClient) WaitForTranscriptionCompleted(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeTranscriptionsMessageCompleted,
	}, timeout)
}

// WaitForTranscriptionsCreated waits for transcriptions to be created
func (c *TranscriptionsClient) WaitForTranscriptionsCreated(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeTranscriptionsCreated,
	}, timeout)
}

// TranscriptionsEventHandler provides default handlers for transcriptions events
type TranscriptionsEventHandler struct {
	OnTranscriptionsCreated           func(*websockets.WebSocketEvent) error
	OnTranscriptionsUpdated           func(*websockets.WebSocketEvent) error
	OnInputAudioBufferCompleted       func(*websockets.WebSocketEvent) error
	OnInputAudioBufferCleared         func(*websockets.WebSocketEvent) error
	OnTranscriptionsMessageUpdate     func(*websockets.TranscriptionsMessageUpdateEvent) error
	OnTranscriptionsMessageCompleted  func(*websockets.WebSocketEvent) error
	OnError                           func(error) error
	OnClosed                          func() error
}

// RegisterHandlers registers all handlers with the client
func (h *TranscriptionsEventHandler) RegisterHandlers(client *TranscriptionsClient) {
	if h.OnTranscriptionsCreated != nil {
		client.OnEvent(websockets.EventTypeTranscriptionsCreated, h.OnTranscriptionsCreated)
	}
	
	if h.OnTranscriptionsUpdated != nil {
		client.OnEvent(websockets.EventTypeTranscriptionsUpdated, h.OnTranscriptionsUpdated)
	}
	
	if h.OnInputAudioBufferCompleted != nil {
		client.OnEvent(websockets.EventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}
	
	if h.OnInputAudioBufferCleared != nil {
		client.OnEvent(websockets.EventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}
	
	if h.OnTranscriptionsMessageUpdate != nil {
		client.OnEvent(websockets.EventTypeTranscriptionsMessageUpdate, func(event *websockets.WebSocketEvent) error {
			var messageEvent websockets.TranscriptionsMessageUpdateEvent
			if err := json.Unmarshal(event.Data, &messageEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal transcriptions message update event: %w", err)
			}
			messageEvent.EventType = event.EventType
			messageEvent.ID = event.ID
			messageEvent.Detail = event.Detail
			return h.OnTranscriptionsMessageUpdate(&messageEvent)
		})
	}
	
	if h.OnTranscriptionsMessageCompleted != nil {
		client.OnEvent(websockets.EventTypeTranscriptionsMessageCompleted, h.OnTranscriptionsMessageCompleted)
	}
	
	if h.OnError != nil {
		client.OnEvent(websockets.EventTypeError, func(event *websockets.WebSocketEvent) error {
			return h.OnError(fmt.Errorf("WebSocket error: %s", string(event.Data)))
		})
	}
	
	if h.OnClosed != nil {
		client.OnEvent(websockets.EventTypeClosed, func(event *websockets.WebSocketEvent) error {
			return h.OnClosed()
		})
	}
}