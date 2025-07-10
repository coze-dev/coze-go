package coze

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Transcriptions returns a new transcriptions client
func (a *websocketAudio) Transcriptions() *TranscriptionsClient {
	return NewTranscriptionsClient(a.baseURL, a.auth, opts...)
}

// TranscriptionsClient handles audio transcriptions WebSocket connections
type TranscriptionsClient struct {
	ws *websocketClient
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
	wsClient := newWebSocketClient(
		&WebSocketClientOption{
			BaseURL: baseURL,
			Path:    "/v1/audio/transcriptions",
			Auth:    auth,
		},
	)

	client := &TranscriptionsClient{
		ws: wsClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *TranscriptionsClient) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *TranscriptionsClient) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *TranscriptionsClient) IsConnected() bool {
	return c.ws.IsConnected()
}

// UpdateTranscriptions updates the transcriptions configuration
func (c *TranscriptionsClient) UpdateTranscriptions(inputAudio *InputAudio) error {
	event := WebSocketTranscriptionsUpdateEvent{
		EventType: WebSocketEventTypeTranscriptionsUpdate,
		Data: &WebSocketTranscriptionsUpdateEventData{
			InputAudio: inputAudio,
		},
	}

	return c.ws.sendEvent(event)
}

// AppendAudioBuffer appends audio data to the input buffer
func (c *TranscriptionsClient) AppendAudioBuffer(audioData []byte) error {
	// Encode audio data to base64
	encoded := base64.StdEncoding.EncodeToString(audioData)

	event := WebSocketInputAudioBufferAppendEvent{
		EventType: WebSocketEventTypeInputAudioBufferAppend,
		Data: &WebSocketInputAudioBufferAppendEventData{
			Delta: encoded,
		},
	}

	return c.ws.sendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *TranscriptionsClient) CompleteAudioBuffer() error {
	event := WebSocketInputAudioBufferCompleteEvent{
		EventType: WebSocketEventTypeInputAudioBufferComplete,
	}

	return c.ws.sendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *TranscriptionsClient) ClearAudioBuffer() error {
	event := WebSocketInputAudioBufferClearEvent{
		EventType: WebSocketEventTypeInputAudioBufferClear,
	}

	return c.ws.sendEvent(event)
}

// OnEvent registers an event handler
func (c *TranscriptionsClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

// WaitForTranscriptionCompleted waits for transcription to complete
func (c *TranscriptionsClient) WaitForTranscriptionCompleted(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeTranscriptionsMessageCompleted,
	}, timeout)
}

// WaitForTranscriptionsCreated waits for transcriptions to be created
func (c *TranscriptionsClient) WaitForTranscriptionsCreated(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeTranscriptionsCreated,
	}, timeout)
}

// TranscriptionsEventHandler provides default handlers for transcriptions events
type TranscriptionsEventHandler struct {
	OnTranscriptionsCreated          func(IWebSocketEvent) error
	OnTranscriptionsUpdated          func(IWebSocketEvent) error
	OnInputAudioBufferCompleted      func(IWebSocketEvent) error
	OnInputAudioBufferCleared        func(IWebSocketEvent) error
	OnTranscriptionsMessageUpdate    func(*WebSocketTranscriptionsMessageUpdateEvent) error
	OnTranscriptionsMessageCompleted func(IWebSocketEvent) error
	OnError                          func(error) error
	OnClosed                         func() error
}

// RegisterHandlers registers all handlers with the client
func (h *TranscriptionsEventHandler) RegisterHandlers(client *TranscriptionsClient) {
	if h.OnTranscriptionsCreated != nil {
		client.OnEvent(WebSocketEventTypeTranscriptionsCreated, h.OnTranscriptionsCreated)
	}

	if h.OnTranscriptionsUpdated != nil {
		client.OnEvent(WebSocketEventTypeTranscriptionsUpdated, h.OnTranscriptionsUpdated)
	}

	if h.OnInputAudioBufferCompleted != nil {
		client.OnEvent(WebSocketEventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}

	if h.OnInputAudioBufferCleared != nil {
		client.OnEvent(WebSocketEventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}

	if h.OnTranscriptionsMessageUpdate != nil {
		client.OnEvent(WebSocketEventTypeTranscriptionsMessageUpdate, func(event IWebSocketEvent) error {
			var messageEvent WebSocketTranscriptionsMessageUpdateEvent
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
		client.OnEvent(WebSocketEventTypeTranscriptionsMessageCompleted, h.OnTranscriptionsMessageCompleted)
	}

	if h.OnError != nil {
		client.OnEvent(WebSocketEventTypeError, func(event IWebSocketEvent) error {
			return h.OnError(fmt.Errorf("WebSocket error: %s", string(event.Data)))
		})
	}

	if h.OnClosed != nil {
		client.OnEvent(WebSocketEventTypeClosed, func(event IWebSocketEvent) error {
			return h.OnClosed()
		})
	}
}
