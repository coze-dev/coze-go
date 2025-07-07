package chat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-go/websockets"
)

// ChatClient handles chat WebSocket connections
type ChatClient struct {
	wsClient *websockets.WebSocketClient
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

// WithInputAudio configures input audio settings
func WithInputAudio(inputAudio *websockets.InputAudio) ChatClientOption {
	return func(c *ChatClient) {
		// This will be sent during connection setup
	}
}

// WithOutputAudio configures output audio settings
func WithOutputAudio(outputAudio *websockets.OutputAudio) ChatClientOption {
	return func(c *ChatClient) {
		// This will be sent during connection setup
	}
}

// NewChatClient creates a new chat WebSocket client
func NewChatClient(baseURL string, auth websockets.Auth, opts ...ChatClientOption) *ChatClient {
	wsClient := websockets.NewWebSocketClient(
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
func (c *ChatClient) UpdateChat(botID string, inputAudio *websockets.InputAudio, outputAudio *websockets.OutputAudio) error {
	event := websockets.ChatUpdateEvent{
		EventType: websockets.EventTypeChatUpdate,
		Data: &websockets.ChatUpdateData{
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
	
	event := websockets.InputAudioBufferAppendEvent{
		EventType: websockets.EventTypeInputAudioBufferAppend,
		Data: &websockets.InputAudioBufferAppendData{
			Delta: encoded,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *ChatClient) CompleteAudioBuffer() error {
	event := websockets.InputAudioBufferCompleteEvent{
		EventType: websockets.EventTypeInputAudioBufferComplete,
	}
	
	return c.wsClient.SendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *ChatClient) ClearAudioBuffer() error {
	event := websockets.InputAudioBufferClearEvent{
		EventType: websockets.EventTypeInputAudioBufferClear,
	}
	
	return c.wsClient.SendEvent(event)
}

// CreateMessage creates a conversation message
func (c *ChatClient) CreateMessage(content string) error {
	event := websockets.ConversationMessageCreateEvent{
		EventType: websockets.EventTypeConversationMessageCreate,
		Data: &websockets.ConversationMessageCreateData{
			Content: content,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// ClearConversation clears the conversation context
func (c *ChatClient) ClearConversation() error {
	event := websockets.WebSocketEvent{
		EventType: websockets.EventTypeConversationClear,
	}
	
	return c.wsClient.SendEvent(event)
}

// SubmitToolOutputs submits tool outputs for a chat
func (c *ChatClient) SubmitToolOutputs(chatID string, toolOutputs []websockets.ToolOutput) error {
	event := websockets.ConversationChatSubmitToolOutputsEvent{
		EventType: websockets.EventTypeConversationChatSubmitToolOutputs,
		Data: &websockets.ConversationChatSubmitToolOutputsData{
			ChatID:      chatID,
			ToolOutputs: toolOutputs,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// CancelChat cancels the current chat
func (c *ChatClient) CancelChat(chatID string) error {
	event := websockets.ConversationChatCancelEvent{
		EventType: websockets.EventTypeConversationChatCancel,
		Data: &websockets.ConversationChatCancelData{
			ChatID: chatID,
		},
	}
	
	return c.wsClient.SendEvent(event)
}

// OnEvent registers an event handler
func (c *ChatClient) OnEvent(eventType websockets.WebSocketEventType, handler websockets.EventHandler) {
	c.wsClient.OnEvent(eventType, handler)
}

// WaitForChatCompleted waits for chat to complete
func (c *ChatClient) WaitForChatCompleted(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeConversationChatCompleted,
		websockets.EventTypeConversationChatFailed,
	}, timeout)
}

// WaitForChatCreated waits for chat to be created
func (c *ChatClient) WaitForChatCreated(timeout time.Duration) (*websockets.WebSocketEvent, error) {
	return c.wsClient.WaitForEvent([]websockets.WebSocketEventType{
		websockets.EventTypeConversationChatCreated,
	}, timeout)
}

// ChatEventHandler provides default handlers for chat events
type ChatEventHandler struct {
	OnChatCreated                       func(*websockets.WebSocketEvent) error
	OnChatUpdated                       func(*websockets.WebSocketEvent) error
	OnConversationChatCreated           func(*websockets.ConversationChatCreatedEvent) error
	OnConversationChatInProgress        func(*websockets.WebSocketEvent) error
	OnConversationMessageDelta          func(*websockets.ConversationMessageDeltaEvent) error
	OnConversationAudioSentenceStart    func(*websockets.WebSocketEvent) error
	OnConversationAudioDelta            func(*websockets.ConversationAudioDeltaEvent) error
	OnConversationMessageCompleted      func(*websockets.WebSocketEvent) error
	OnConversationAudioCompleted        func(*websockets.WebSocketEvent) error
	OnConversationChatCompleted         func(*websockets.ConversationChatCompletedEvent) error
	OnConversationChatFailed            func(*websockets.WebSocketEvent) error
	OnInputAudioBufferCompleted         func(*websockets.WebSocketEvent) error
	OnInputAudioBufferCleared           func(*websockets.WebSocketEvent) error
	OnConversationCleared               func(*websockets.WebSocketEvent) error
	OnConversationChatCanceled          func(*websockets.ConversationChatCanceledEvent) error
	OnConversationAudioTranscriptUpdate   func(*websockets.WebSocketEvent) error
	OnConversationAudioTranscriptCompleted func(*websockets.WebSocketEvent) error
	OnConversationChatRequiresAction      func(*websockets.ConversationChatRequiresActionEvent) error
	OnInputAudioBufferSpeechStarted       func(*websockets.WebSocketEvent) error
	OnInputAudioBufferSpeechStopped       func(*websockets.WebSocketEvent) error
	OnError                               func(error) error
	OnClosed                              func() error
}

// RegisterHandlers registers all handlers with the client
func (h *ChatEventHandler) RegisterHandlers(client *ChatClient) {
	if h.OnChatCreated != nil {
		client.OnEvent(websockets.EventTypeChatCreated, h.OnChatCreated)
	}
	
	if h.OnChatUpdated != nil {
		client.OnEvent(websockets.EventTypeChatUpdated, h.OnChatUpdated)
	}
	
	if h.OnConversationChatCreated != nil {
		client.OnEvent(websockets.EventTypeConversationChatCreated, func(event *websockets.WebSocketEvent) error {
			var chatEvent websockets.ConversationChatCreatedEvent
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
		client.OnEvent(websockets.EventTypeConversationChatInProgress, h.OnConversationChatInProgress)
	}
	
	if h.OnConversationMessageDelta != nil {
		client.OnEvent(websockets.EventTypeConversationMessageDelta, func(event *websockets.WebSocketEvent) error {
			var deltaEvent websockets.ConversationMessageDeltaEvent
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
		client.OnEvent(websockets.EventTypeConversationAudioSentenceStart, h.OnConversationAudioSentenceStart)
	}
	
	if h.OnConversationAudioDelta != nil {
		client.OnEvent(websockets.EventTypeConversationAudioDelta, func(event *websockets.WebSocketEvent) error {
			var audioEvent websockets.ConversationAudioDeltaEvent
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
		client.OnEvent(websockets.EventTypeConversationMessageCompleted, h.OnConversationMessageCompleted)
	}
	
	if h.OnConversationAudioCompleted != nil {
		client.OnEvent(websockets.EventTypeConversationAudioCompleted, h.OnConversationAudioCompleted)
	}
	
	if h.OnConversationChatCompleted != nil {
		client.OnEvent(websockets.EventTypeConversationChatCompleted, func(event *websockets.WebSocketEvent) error {
			var completedEvent websockets.ConversationChatCompletedEvent
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
		client.OnEvent(websockets.EventTypeConversationChatFailed, h.OnConversationChatFailed)
	}
	
	if h.OnInputAudioBufferCompleted != nil {
		client.OnEvent(websockets.EventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}
	
	if h.OnInputAudioBufferCleared != nil {
		client.OnEvent(websockets.EventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}
	
	if h.OnConversationCleared != nil {
		client.OnEvent(websockets.EventTypeConversationCleared, h.OnConversationCleared)
	}
	
	if h.OnConversationChatCanceled != nil {
		client.OnEvent(websockets.EventTypeConversationChatCanceled, func(event *websockets.WebSocketEvent) error {
			var canceledEvent websockets.ConversationChatCanceledEvent
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
		client.OnEvent(websockets.EventTypeConversationAudioTranscriptUpdate, h.OnConversationAudioTranscriptUpdate)
	}
	
	if h.OnConversationAudioTranscriptCompleted != nil {
		client.OnEvent(websockets.EventTypeConversationAudioTranscriptCompleted, h.OnConversationAudioTranscriptCompleted)
	}
	
	if h.OnConversationChatRequiresAction != nil {
		client.OnEvent(websockets.EventTypeConversationChatRequiresAction, func(event *websockets.WebSocketEvent) error {
			var actionEvent websockets.ConversationChatRequiresActionEvent
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
		client.OnEvent(websockets.EventTypeInputAudioBufferSpeechStarted, h.OnInputAudioBufferSpeechStarted)
	}
	
	if h.OnInputAudioBufferSpeechStopped != nil {
		client.OnEvent(websockets.EventTypeInputAudioBufferSpeechStopped, h.OnInputAudioBufferSpeechStopped)
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