package coze

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// SpeechClient handles audio speech WebSocket connections
type SpeechClient struct {
	ws *websocketClient
}

// NewSpeechClient creates a new speech WebSocket client
func NewSpeechClient(baseURL string, auth Auth, opts ...SpeechClientOption) *SpeechClient {
	wsClient := newWebSocketClient(
		&WebSocketClientOption{
			BaseURL: baseURL,
			Path:    "/v1/audio/speech",
			Auth:    auth,
		},
	)

	client := &SpeechClient{
		ws: wsClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *SpeechClient) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *SpeechClient) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *SpeechClient) IsConnected() bool {
	return c.ws.IsConnected()
}

// UpdateSpeech updates the speech configuration
func (c *SpeechClient) UpdateSpeech(outputAudio *OutputAudio) error {
	event := SpeechUpdateEvent{
		EventType: EventTypeSpeechUpdate,
		Data: &SpeechUpdateData{
			OutputAudio: outputAudio,
		},
	}

	return c.ws.sendEvent(event)
}

// AppendTextBuffer appends text to the input buffer
func (c *SpeechClient) AppendTextBuffer(text string) error {
	event := InputTextBufferAppendEvent{
		EventType: EventTypeInputTextBufferAppend,
		Data: &InputTextBufferAppendData{
			Delta: text,
		},
	}

	return c.ws.sendEvent(event)
}

// CompleteTextBuffer completes the text buffer input
func (c *SpeechClient) CompleteTextBuffer() error {
	event := InputTextBufferCompleteEvent{
		EventType: EventTypeInputTextBufferComplete,
	}

	return c.ws.sendEvent(event)
}

// OnEvent registers an event handler
func (c *SpeechClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

// WaitForSpeechAudioCompleted waits for speech audio to complete
func (c *SpeechClient) WaitForSpeechAudioCompleted(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		EventTypeSpeechAudioCompleted,
	}, timeout)
}

// WaitForSpeechCreated waits for speech to be created
func (c *SpeechClient) WaitForSpeechCreated(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		EventTypeSpeechCreated,
	}, timeout)
}

// SpeechEventHandler provides default handlers for speech events
type SpeechEventHandler struct {
	OnSpeechCreated            func(*SpeechCreatedEvent) error
	OnSpeechUpdated            func(IWebSocketEvent) error
	OnInputTextBufferCompleted func(IWebSocketEvent) error
	OnSpeechAudioUpdate        func(*SpeechAudioUpdateEvent) error
	OnSpeechAudioCompleted     func(*SpeechAudioCompletedEvent) error
	OnError                    func(error) error
	OnClosed                   func() error
}

// RegisterHandlers registers all handlers with the client
func (h *SpeechEventHandler) RegisterHandlers(client *SpeechClient) {
	if h.OnSpeechCreated != nil {
		client.OnEvent(EventTypeSpeechCreated, func(event IWebSocketEvent) error {
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
		client.OnEvent(EventTypeSpeechAudioUpdate, func(event IWebSocketEvent) error {
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
		client.OnEvent(EventTypeSpeechAudioCompleted, func(event IWebSocketEvent) error {
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
		client.OnEvent(EventTypeError, func(event IWebSocketEvent) error {
			return h.OnError(fmt.Errorf("WebSocket error: %s", string(event.Data)))
		})
	}

	if h.OnClosed != nil {
		client.OnEvent(EventTypeClosed, func(event IWebSocketEvent) error {
			return h.OnClosed()
		})
	}
}

// GetAudioFromDelta extracts audio bytes from a delta string
func GetAudioFromDelta(delta string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(delta)
}
