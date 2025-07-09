package websockets

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// SpeechClient handles audio speech WebSocket connections
type SpeechClient struct {
	wsClient *WebSocketClient
}

// SpeechClientOption configures the speech client
type SpeechClientOption func(*SpeechClient)

// WithOutputAudio configures output audio settings
func WithOutputAudio(outputAudio *OutputAudio) SpeechClientOption {
	return func(c *SpeechClient) {
		// This will be sent during connection setup
	}
}

// NewSpeechClient creates a new speech WebSocket client
func NewSpeechClient(baseURL string, auth Auth, opts ...SpeechClientOption) *SpeechClient {
	wsClient := NewWebSocketClient(
		&WebSocketClientOption{
			BaseURL: baseURL,
			Path:    "/v1/audio/speech",
			Auth:    auth,
		},
	)

	client := &SpeechClient{
		wsClient: wsClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *SpeechClient) Connect() error {
	return c.wsClient.Connect()
}

// Close closes the WebSocket connection
func (c *SpeechClient) Close() error {
	return c.wsClient.Close()
}

// IsConnected returns whether the client is connected
func (c *SpeechClient) IsConnected() bool {
	return c.wsClient.IsConnected()
}

// UpdateSpeech updates the speech configuration
func (c *SpeechClient) UpdateSpeech(outputAudio *OutputAudio) error {
	event := SpeechUpdateEvent{
		EventType: EventTypeSpeechUpdate,
		Data: &SpeechUpdateData{
			OutputAudio: outputAudio,
		},
	}

	return c.wsClient.SendEvent(event)
}

// AppendTextBuffer appends text to the input buffer
func (c *SpeechClient) AppendTextBuffer(text string) error {
	event := InputTextBufferAppendEvent{
		EventType: EventTypeInputTextBufferAppend,
		Data: &InputTextBufferAppendData{
			Delta: text,
		},
	}

	return c.wsClient.SendEvent(event)
}

// CompleteTextBuffer completes the text buffer input
func (c *SpeechClient) CompleteTextBuffer() error {
	event := InputTextBufferCompleteEvent{
		EventType: EventTypeInputTextBufferComplete,
	}

	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *SpeechClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForSpeechAudioCompleted waits for speech audio to complete
func (c *SpeechClient) WaitForSpeechAudioCompleted(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeSpeechAudioCompleted,
	}, timeout)
}

// WaitForSpeechCreated waits for speech to be created
func (c *SpeechClient) WaitForSpeechCreated(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeSpeechCreated,
	}, timeout)
}

// SpeechEventHandler provides default handlers for speech events
type SpeechEventHandler struct {
	OnSpeechCreated            func(*SpeechCreatedEvent) error
	OnSpeechUpdated            func(*WebSocketEvent) error
	OnInputTextBufferCompleted func(*WebSocketEvent) error
	OnSpeechAudioUpdate        func(*SpeechAudioUpdateEvent) error
	OnSpeechAudioCompleted     func(*SpeechAudioCompletedEvent) error
	OnError                    func(error) error
	OnClosed                   func() error
}

// RegisterHandlers registers all handlers with the client
func (h *SpeechEventHandler) RegisterHandlers(client *SpeechClient) {
	if h.OnSpeechCreated != nil {
		client.OnEvent(EventTypeSpeechCreated, func(event *WebSocketEvent) error {
			var speechEvent SpeechCreatedEvent
			if err := json.Unmarshal(event.Data, &speechEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal speech created event: %w", err)
			}
			speechEvent.EventType = event.EventType
			speechEvent.ID = event.ID
			speechEvent.Detail = event.Detail
			return h.OnSpeechCreated(&speechEvent)
		})
	}

	if h.OnSpeechUpdated != nil {
		client.OnEvent(EventTypeSpeechUpdated, h.OnSpeechUpdated)
	}

	if h.OnInputTextBufferCompleted != nil {
		client.OnEvent(EventTypeInputTextBufferCompleted, h.OnInputTextBufferCompleted)
	}

	if h.OnSpeechAudioUpdate != nil {
		client.OnEvent(EventTypeSpeechAudioUpdate, func(event *WebSocketEvent) error {
			var audioEvent SpeechAudioUpdateEvent
			if err := json.Unmarshal(event.Data, &audioEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal speech audio update event: %w", err)
			}
			audioEvent.EventType = event.EventType
			audioEvent.ID = event.ID
			audioEvent.Detail = event.Detail
			return h.OnSpeechAudioUpdate(&audioEvent)
		})
	}

	if h.OnSpeechAudioCompleted != nil {
		client.OnEvent(EventTypeSpeechAudioCompleted, func(event *WebSocketEvent) error {
			var completedEvent SpeechAudioCompletedEvent
			if err := json.Unmarshal(event.Data, &completedEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal speech audio completed event: %w", err)
			}
			completedEvent.EventType = event.EventType
			completedEvent.ID = event.ID
			completedEvent.Detail = event.Detail
			return h.OnSpeechAudioCompleted(&completedEvent)
		})
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

// GetAudioFromDelta extracts audio bytes from a delta string
func GetAudioFromDelta(delta string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(delta)
}
