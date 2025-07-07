package websockets

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// ChatClient handles chat WebSocket connections
type ChatClient struct {
	wsClient *WebSocketClient
	botID    string
}

// ChatClientOption configures the chat client
type ChatClientOption func(*ChatClient)

// WithBotID sets the bot ID for the chat
func WithBotID(botID string) ChatClientOption {
	return func(c *ChatClient) {
		c.botID = botID
	}
}

// WithChatInputAudio configures input audio settings
func WithChatInputAudio(inputAudio *InputAudio) ChatClientOption {
	return func(c *ChatClient) {
		// This will be sent during connection setup
	}
}

// WithChatOutputAudio configures output audio settings
func WithChatOutputAudio(outputAudio *OutputAudio) ChatClientOption {
	return func(c *ChatClient) {
		// This will be sent during connection setup
	}
}

// NewChatClient creates a new chat WebSocket client
func NewChatClient(baseURL string, auth Auth, opts ...ChatClientOption) *ChatClient {
	wsClient := NewWebSocketClient(
		baseURL,
		"v1/chat",
		auth,
	)
	
	client := &ChatClient{
		wsClient: wsClient,
	}
	
	for _, opt := range opts {
		opt(client)
	}
	
	return client
}

// Connect establishes the WebSocket connection
func (c *ChatClient) Connect() error {
	return c.wsClient.Connect()
}

// Close closes the WebSocket connection
func (c *ChatClient) Close() error {
	return c.wsClient.Close()
}

// IsConnected returns whether the client is connected
func (c *ChatClient) IsConnected() bool {
	return c.wsClient.IsConnected()
}

// UpdateChat updates the chat configuration
func (c *ChatClient) UpdateChat(botID string, inputAudio *InputAudio, outputAudio *OutputAudio) error {
	event := ChatUpdateEvent{
		EventType: EventTypeChatUpdate,
		Data: &ChatUpdateData{
			BotID:       botID,
			InputAudio:  inputAudio,
			OutputAudio: outputAudio,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// AppendAudioBuffer appends audio data to the input buffer
func (c *ChatClient) AppendAudioBuffer(audioData []byte) error {
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
func (c *ChatClient) CompleteAudioBuffer() error {
	event := InputAudioBufferCompleteEvent{
		EventType: EventTypeInputAudioBufferComplete,
	}
	
	return c.wsClient.SendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *ChatClient) ClearAudioBuffer() error {
	event := InputAudioBufferClearEvent{
		EventType: EventTypeInputAudioBufferClear,
	}
	
	return c.wsClient.SendEvent(event)
}

// CreateMessage creates a conversation message
func (c *ChatClient) CreateMessage(content string) error {
	event := ConversationMessageCreateEvent{
		EventType: EventTypeConversationMessageCreate,
		Data: &ConversationMessageCreateData{
			Content: content,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// ClearConversation clears the conversation context
func (c *ChatClient) ClearConversation() error {
	event := WebSocketEvent{
		EventType: EventTypeConversationClear,
	}
	
	return c.wsClient.SendEvent(event)
}

// SubmitToolOutputs submits tool outputs for a chat
func (c *ChatClient) SubmitToolOutputs(chatID string, toolOutputs []ToolOutput) error {
	event := ConversationChatSubmitToolOutputsEvent{
		EventType: EventTypeConversationChatSubmitToolOutputs,
		Data: &ConversationChatSubmitToolOutputsData{
			ChatID:      chatID,
			ToolOutputs: toolOutputs,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// CancelChat cancels the current chat
func (c *ChatClient) CancelChat(chatID string) error {
	event := ConversationChatCancelEvent{
		EventType: EventTypeConversationChatCancel,
		Data: &ConversationChatCancelData{
			ChatID: chatID,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *ChatClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForChatCompleted waits for chat to complete
func (c *ChatClient) WaitForChatCompleted(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeConversationChatCompleted,
		EventTypeConversationChatFailed,
	}, timeout)
}

// WaitForChatCreated waits for chat to be created
func (c *ChatClient) WaitForChatCreated(timeout time.Duration) (*WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]WebSocketEventType{
		EventTypeConversationChatCreated,
	}, timeout)
}

// ChatEventHandler provides default handlers for chat events
type ChatEventHandler struct {
	OnChatCreated                       func(*WebSocketEvent) error
	OnChatUpdated                       func(*WebSocketEvent) error
	OnConversationChatCreated           func(*ConversationChatCreatedEvent) error
	OnConversationChatInProgress        func(*WebSocketEvent) error
	OnConversationMessageDelta          func(*ConversationMessageDeltaEvent) error
	OnConversationAudioSentenceStart    func(*WebSocketEvent) error
	OnConversationAudioDelta            func(*ConversationAudioDeltaEvent) error
	OnConversationMessageCompleted      func(*WebSocketEvent) error
	OnConversationAudioCompleted        func(*WebSocketEvent) error
	OnConversationChatCompleted         func(*ConversationChatCompletedEvent) error
	OnConversationChatFailed            func(*WebSocketEvent) error
	OnInputAudioBufferCompleted         func(*WebSocketEvent) error
	OnInputAudioBufferCleared           func(*WebSocketEvent) error
	OnConversationCleared               func(*WebSocketEvent) error
	OnConversationChatCanceled          func(*ConversationChatCanceledEvent) error
	OnConversationAudioTranscriptUpdate   func(*WebSocketEvent) error
	OnConversationAudioTranscriptCompleted func(*WebSocketEvent) error
	OnConversationChatRequiresAction      func(*ConversationChatRequiresActionEvent) error
	OnInputAudioBufferSpeechStarted       func(*WebSocketEvent) error
	OnInputAudioBufferSpeechStopped       func(*WebSocketEvent) error
	OnError                               func(error) error
	OnClosed                              func() error
}

// RegisterHandlers registers all handlers with the client
func (h *ChatEventHandler) RegisterHandlers(client *ChatClient) {
	if h.OnChatCreated != nil {
		client.OnEvent(EventTypeChatCreated, h.OnChatCreated)
	}
	
	if h.OnChatUpdated != nil {
		client.OnEvent(EventTypeChatUpdated, h.OnChatUpdated)
	}
	
	if h.OnConversationChatCreated != nil {
		client.OnEvent(EventTypeConversationChatCreated, func(event *WebSocketEvent) error {
			var chatEvent ConversationChatCreatedEvent
			if err := json.Unmarshal(event.Data, &chatEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation chat created event: %w", err)
			}
			chatEvent.EventType = event.EventType
			chatEvent.ID = event.ID
			chatEvent.Detail = event.Detail
			return h.OnConversationChatCreated(&chatEvent)
		})
	}
	
	if h.OnConversationChatInProgress != nil {
		client.OnEvent(EventTypeConversationChatInProgress, h.OnConversationChatInProgress)
	}
	
	if h.OnConversationMessageDelta != nil {
		client.OnEvent(EventTypeConversationMessageDelta, func(event *WebSocketEvent) error {
			var deltaEvent ConversationMessageDeltaEvent
			if err := json.Unmarshal(event.Data, &deltaEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation message delta event: %w", err)
			}
			deltaEvent.EventType = event.EventType
			deltaEvent.ID = event.ID
			deltaEvent.Detail = event.Detail
			return h.OnConversationMessageDelta(&deltaEvent)
		})
	}
	
	if h.OnConversationAudioSentenceStart != nil {
		client.OnEvent(EventTypeConversationAudioSentenceStart, h.OnConversationAudioSentenceStart)
	}
	
	if h.OnConversationAudioDelta != nil {
		client.OnEvent(EventTypeConversationAudioDelta, func(event *WebSocketEvent) error {
			var audioEvent ConversationAudioDeltaEvent
			if err := json.Unmarshal(event.Data, &audioEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation audio delta event: %w", err)
			}
			audioEvent.EventType = event.EventType
			audioEvent.ID = event.ID
			audioEvent.Detail = event.Detail
			return h.OnConversationAudioDelta(&audioEvent)
		})
	}
	
	if h.OnConversationMessageCompleted != nil {
		client.OnEvent(EventTypeConversationMessageCompleted, h.OnConversationMessageCompleted)
	}
	
	if h.OnConversationAudioCompleted != nil {
		client.OnEvent(EventTypeConversationAudioCompleted, h.OnConversationAudioCompleted)
	}
	
	if h.OnConversationChatCompleted != nil {
		client.OnEvent(EventTypeConversationChatCompleted, func(event *WebSocketEvent) error {
			var completedEvent ConversationChatCompletedEvent
			if err := json.Unmarshal(event.Data, &completedEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation chat completed event: %w", err)
			}
			completedEvent.EventType = event.EventType
			completedEvent.ID = event.ID
			completedEvent.Detail = event.Detail
			return h.OnConversationChatCompleted(&completedEvent)
		})
	}
	
	if h.OnConversationChatFailed != nil {
		client.OnEvent(EventTypeConversationChatFailed, h.OnConversationChatFailed)
	}
	
	if h.OnInputAudioBufferCompleted != nil {
		client.OnEvent(EventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}
	
	if h.OnInputAudioBufferCleared != nil {
		client.OnEvent(EventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}
	
	if h.OnConversationCleared != nil {
		client.OnEvent(EventTypeConversationCleared, h.OnConversationCleared)
	}
	
	if h.OnConversationChatCanceled != nil {
		client.OnEvent(EventTypeConversationChatCanceled, func(event *WebSocketEvent) error {
			var canceledEvent ConversationChatCanceledEvent
			if err := json.Unmarshal(event.Data, &canceledEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation chat canceled event: %w", err)
			}
			canceledEvent.EventType = event.EventType
			canceledEvent.ID = event.ID
			canceledEvent.Detail = event.Detail
			return h.OnConversationChatCanceled(&canceledEvent)
		})
	}
	
	if h.OnConversationAudioTranscriptUpdate != nil {
		client.OnEvent(EventTypeConversationAudioTranscriptUpdate, h.OnConversationAudioTranscriptUpdate)
	}
	
	if h.OnConversationAudioTranscriptCompleted != nil {
		client.OnEvent(EventTypeConversationAudioTranscriptCompleted, h.OnConversationAudioTranscriptCompleted)
	}
	
	if h.OnConversationChatRequiresAction != nil {
		client.OnEvent(EventTypeConversationChatRequiresAction, func(event *WebSocketEvent) error {
			var actionEvent ConversationChatRequiresActionEvent
			if err := json.Unmarshal(event.Data, &actionEvent.Data); err != nil {
				return fmt.Errorf("failed to unmarshal conversation chat requires action event: %w", err)
			}
			actionEvent.EventType = event.EventType
			actionEvent.ID = event.ID
			actionEvent.Detail = event.Detail
			return h.OnConversationChatRequiresAction(&actionEvent)
		})
	}
	
	if h.OnInputAudioBufferSpeechStarted != nil {
		client.OnEvent(EventTypeInputAudioBufferSpeechStarted, h.OnInputAudioBufferSpeechStarted)
	}
	
	if h.OnInputAudioBufferSpeechStopped != nil {
		client.OnEvent(EventTypeInputAudioBufferSpeechStopped, h.OnInputAudioBufferSpeechStopped)
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