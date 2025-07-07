package audio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-go/websockets"
)

// SpeechClient handles audio speech WebSocket connections
type SpeechClient struct {
	wsClient *websockets.WebSocketClient
}

// SpeechClientOption configures the speech client
type SpeechClientOption func(*SpeechClient)

// WithOutputAudio configures output audio settings
func WithOutputAudio(outputAudio *websockets.OutputAudio) SpeechClientOption {
	return func(c *SpeechClient) {
		// This will be sent during connection setup
	}
}

// NewSpeechClient creates a new speech WebSocket client
func NewSpeechClient(baseURL string, auth websockets.Auth, opts ...SpeechClientOption) *SpeechClient {
	wsClient := websockets.NewWebSocketClient(
		baseURL,
		"v1/audio/speech",
		auth,
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
func (c *SpeechClient) UpdateSpeech(outputAudio *websockets.OutputAudio) error {
	event := websockets.SpeechUpdateEvent{
		EventType: websockets.EventTypeSpeechUpdate,
		Data: &websockets.SpeechUpdateData{
			OutputAudio: outputAudio,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// AppendTextBuffer appends text to the input buffer
func (c *SpeechClient) AppendTextBuffer(text string) error {
	event := websockets.InputTextBufferAppendEvent{
		EventType: websockets.EventTypeInputTextBufferAppend,
		Data: &websockets.InputTextBufferAppendData{
			Delta: text,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// CompleteTextBuffer completes the text buffer input
func (c *SpeechClient) CompleteTextBuffer() error {
	event := websockets.InputTextBufferCompleteEvent{
		EventType: websockets.EventTypeInputTextBufferComplete,
	}
	
	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *SpeechClient) OnEvent(eventType websockets.WebSocketEventType, handler websockets.EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForSpeechAudioCompleted waits for speech audio to complete
func (c *SpeechClient) WaitForSpeechAudioCompleted(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeSpeechAudioCompleted,
	}, timeout)
}

// WaitForSpeechCreated waits for speech to be created
func (c *SpeechClient) WaitForSpeechCreated(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeSpeechCreated,
	}, timeout)
}

// SpeechEventHandler provides default handlers for speech events
type SpeechEventHandler struct {
	OnSpeechCreated           func(*websockets.SpeechCreatedEvent) error
	OnSpeechUpdated           func(*websockets.WebSocketEvent) error
	OnInputTextBufferCompleted func(*websockets.WebSocketEvent) error
	OnSpeechAudioUpdate       func(*websockets.SpeechAudioUpdateEvent) error
	OnSpeechAudioCompleted    func(*websockets.SpeechAudioCompletedEvent) error
	OnError                   func(error) error
	OnClosed                  func() error
}

// RegisterHandlers registers all handlers with the client
func (h *SpeechEventHandler) RegisterHandlers(client *SpeechClient) {
	if h.OnSpeechCreated != nil {
		client.OnEvent(websockets.EventTypeSpeechCreated, func(event *websockets.WebSocketEvent) error {
			var speechEvent websockets.SpeechCreatedEvent
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
		client.OnEvent(websockets.EventTypeSpeechUpdated, h.OnSpeechUpdated)
	}
	
	if h.OnInputTextBufferCompleted != nil {
		client.OnEvent(websockets.EventTypeInputTextBufferCompleted, h.OnInputTextBufferCompleted)
	}
	
	if h.OnSpeechAudioUpdate != nil {
		client.OnEvent(websockets.EventTypeSpeechAudioUpdate, func(event *websockets.WebSocketEvent) error {
			var audioEvent websockets.SpeechAudioUpdateEvent
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
		client.OnEvent(websockets.EventTypeSpeechAudioCompleted, func(event *websockets.WebSocketEvent) error {
			var completedEvent websockets.SpeechAudioCompletedEvent
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

// GetAudioFromDelta extracts audio bytes from a delta string
func GetAudioFromDelta(delta string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(delta)
}

// SpeechClient handles audio speech WebSocket connections
type SpeechClient struct {
	*websockets.WebSocketClient
}

// SpeechClientOption configures the speech client
type SpeechClientOption func(*SpeechClient)

// WithOutputAudio configures output audio settings
func WithOutputAudio(outputAudio *websockets.OutputAudio) SpeechClientOption {
	return func(c *SpeechClient) {
		// This will be sent during connection setup
	}
}

// NewSpeechClient creates a new speech WebSocket client
func NewSpeechClient(baseURL string, auth websockets.Auth, opts ...SpeechClientOption) *SpeechClient {
	wsClient := websockets.NewWebSocketClient(
		baseURL,
		"v1/audio/speech",
		auth,
	)
	
	client := &SpeechClient{
		WebSocketClient: wsClient,
	}
	
	for _, opt := range opts {
		opt(client)
	}
	
	return client
}

// UpdateSpeech updates the speech configuration
func (c *SpeechClient) UpdateSpeech(outputAudio *websockets.OutputAudio) error {
	event := websockets.SpeechUpdateEvent{
		EventType: websockets.EventTypeSpeechUpdate,
		Data: &websockets.SpeechUpdateData{
			OutputAudio: outputAudio,
		},
	}
	
	return c.SendEvent(event)
}

// AppendTextBuffer appends text to the input buffer
func (c *SpeechClient) AppendTextBuffer(text string) error {
	event := websockets.InputTextBufferAppendEvent{
		EventType: websockets.EventTypeInputTextBufferAppend,
		Data: &websockets.InputTextBufferAppendData{
			Delta: text,
		},
	}
	
	return c.SendEvent(event)
}

// CompleteTextBuffer completes the text buffer input
func (c *SpeechClient) CompleteTextBuffer() error {
	event := websockets.InputTextBufferCompleteEvent{
		EventType: websockets.EventTypeInputTextBufferComplete,
	}
	
	return c.SendEvent(event)
}

// WaitForSpeechAudioCompleted waits for speech audio to complete
func (c *SpeechClient) WaitForSpeechAudioCompleted(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeSpeechAudioCompleted,
	}, timeout)
}

// WaitForSpeechCreated waits for speech to be created
func (c *SpeechClient) WaitForSpeechCreated(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeSpeechCreated,
	}, timeout)
}

// SpeechEventHandler provides default handlers for speech events
type SpeechEventHandler struct {
	OnSpeechCreated           func(*websockets.SpeechCreatedEvent) error
	OnSpeechUpdated           func(*websockets.WebSocketEvent) error
	OnInputTextBufferCompleted func(*websockets.WebSocketEvent) error
	OnSpeechAudioUpdate       func(*websockets.SpeechAudioUpdateEvent) error
	OnSpeechAudioCompleted    func(*websockets.SpeechAudioCompletedEvent) error
	OnError                   func(error) error
	OnClosed                  func() error
}

// RegisterHandlers registers all handlers with the client
func (h *SpeechEventHandler) RegisterHandlers(client *SpeechClient) {
	if h.OnSpeechCreated != nil {
		client.OnEvent(websockets.EventTypeSpeechCreated, func(event *websockets.WebSocketEvent) error {
			var speechEvent websockets.SpeechCreatedEvent
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
		client.OnEvent(websockets.EventTypeSpeechUpdated, h.OnSpeechUpdated)
	}
	
	if h.OnInputTextBufferCompleted != nil {
		client.OnEvent(websockets.EventTypeInputTextBufferCompleted, h.OnInputTextBufferCompleted)
	}
	
	if h.OnSpeechAudioUpdate != nil {
		client.OnEvent(websockets.EventTypeSpeechAudioUpdate, func(event *websockets.WebSocketEvent) error {
			var audioEvent websockets.SpeechAudioUpdateEvent
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
		client.OnEvent(websockets.EventTypeSpeechAudioCompleted, func(event *websockets.WebSocketEvent) error {
			var completedEvent websockets.SpeechAudioCompletedEvent
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

// GetAudioFromDelta extracts audio bytes from a delta string
func GetAudioFromDelta(delta string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(delta)
}