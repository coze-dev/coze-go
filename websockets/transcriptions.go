package websockets

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// TranscriptionsClient handles audio transcriptions WebSocket connections
type TranscriptionsClient struct {
	wsClient *WebSocketClient
}

// TranscriptionsClientOption configures the transcriptions client
type TranscriptionsClientOption func(*TranscriptionsClient)

// WithInputAudio configures input audio settings
func WithInputAudio(inputAudio *InputAudio) TranscriptionsClientOption {
	return func(c *TranscriptionsClient) {
		// This will be sent during connection setup
	}
}

// NewTranscriptionsClient creates a new transcriptions WebSocket client
func NewTranscriptionsClient(baseURL string, auth Auth, opts ...TranscriptionsClientOption) *TranscriptionsClient {
	wsClient := NewWebSocketClient(
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
func (c *TranscriptionsClient) UpdateTranscriptions(inputAudio *InputAudio) error {
	event := TranscriptionsUpdateEvent{
		EventType: EventTypeTranscriptionsUpdate,
		Data: &TranscriptionsUpdateData{
			InputAudio: inputAudio,
		},
	}

	return c.wsClient.SendEvent(event)
}

// AppendAudioBuffer appends audio data to the input buffer
func (c *TranscriptionsClient) AppendAudioBuffer(audioData []byte) error {
	// Encode audio data to base64
	encoded := base64.StdEncoding.EncodeToString(audioData)

	event := InputAudioBufferAppendEvent{
		EventType: EventTypeInputAudioBufferAppend,
		Data: &InputAudioBufferAppendData{
			Delta: encoded,
		},
	}

	return c.wsClient.SendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *TranscriptionsClient) CompleteAudioBuffer() error {
	event := InputAudioBufferCompleteEvent{
		EventType: EventTypeInputAudioBufferComplete,
	}

	return c.wsClient.SendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *TranscriptionsClient) ClearAudioBuffer() error {
	event := InputAudioBufferClearEvent{
		EventType: EventTypeInputAudioBufferClear,
	}

	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *TranscriptionsClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForTranscriptionCompleted waits for transcription to complete
func (c *TranscriptionsClient) WaitForTranscriptionCompleted(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeTranscriptionsMessageCompleted,
	}, timeout)
}

// WaitForTranscriptionsCreated waits for transcriptions to be created
func (c *TranscriptionsClient) WaitForTranscriptionsCreated(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeTranscriptionsCreated,
	}, timeout)
}

// TranscriptionsEventHandler provides default handlers for transcriptions events
type TranscriptionsEventHandler struct {
	OnTranscriptionsCreated          func(*WebSocketEvent) error
	OnTranscriptionsUpdated          func(*WebSocketEvent) error
	OnInputAudioBufferCompleted      func(*WebSocketEvent) error
	OnInputAudioBufferCleared        func(*WebSocketEvent) error
	OnTranscriptionsMessageUpdate    func(*TranscriptionsMessageUpdateEvent) error
	OnTranscriptionsMessageCompleted func(*WebSocketEvent) error
	OnError                          func(error) error
	OnClosed                         func() error
}

// RegisterHandlers registers all handlers with the client
func (h *TranscriptionsEventHandler) RegisterHandlers(client *TranscriptionsClient) {
	if h.OnTranscriptionsCreated != nil {
		client.OnEvent(EventTypeTranscriptionsCreated, h.OnTranscriptionsCreated)
	}

	if h.OnTranscriptionsUpdated != nil {
		client.OnEvent(EventTypeTranscriptionsUpdated, h.OnTranscriptionsUpdated)
	}

	if h.OnInputAudioBufferCompleted != nil {
		client.OnEvent(EventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}

	if h.OnInputAudioBufferCleared != nil {
		client.OnEvent(EventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}

	if h.OnTranscriptionsMessageUpdate != nil {
		client.OnEvent(EventTypeTranscriptionsMessageUpdate, func(event *WebSocketEvent) error {
			var messageEvent TranscriptionsMessageUpdateEvent
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
		client.OnEvent(EventTypeTranscriptionsMessageCompleted, h.OnTranscriptionsMessageCompleted)
	}

	if h.OnError != nil {
		client.OnEvent(EventTypeError, func(event *WebSocketEvent) error {
			return h.OnError(fmt.Errorf("WebSocket error: %s", string(event.Data)))
		})
	}

	if h.OnClosed != nil {
		client.OnEvent(EventTypeClosed, func(event *WebSocketEvent) error {
			return h.OnClosed()
		})
	}
}
