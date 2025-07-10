package coze

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

func (c *websocketChatBuilder) Create(ctx context.Context, req *CreateWebsocketChatReq) *websocketsChat {
	return newWebsocketChatClient(ctx, c.core, req)
}

type CreateWebsocketChatReq struct {
	// BotID is the ID of the bot.
	BotID string `json:"bot_id"`
}

type WebsocketXXX interface {
	Connect() error
	Close() error
	IsConnected() bool
	OnEvents(handlers map[WebSocketEventType]EventHandler)
	OnEvent(eventType WebSocketEventType, handler EventHandler)
}

type websocketsChat struct {
	ctx  context.Context
	core *core
	req  *CreateWebsocketChatReq

	ws *websocketClient
}

func newWebsocketChatClient(ctx context.Context, core *core, req *CreateWebsocketChatReq) *websocketsChat {
	ws := newWebSocketClient(
		&WebSocketClientOption{
			ctx:  ctx,
			core: core,
			path: "/v1/chat",
		},
	)

	client := &websocketsChat{
		ws: ws,
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *websocketsChat) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *websocketsChat) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *websocketsChat) IsConnected() bool {
	return c.ws.IsConnected()
}

// UpdateChat updates the chat configuration
func (c *websocketsChat) UpdateChat(botID string, inputAudio *InputAudio, outputAudio *OutputAudio) error {
	event := WebSocketChatUpdateEvent{
		EventType: EventTypeChatUpdate,
		Data: &WebSocketChatUpdateEventData{
			BotID:       botID,
			InputAudio:  inputAudio,
			OutputAudio: outputAudio,
		},
	}

	return c.ws.sendEvent(event)
}

// AppendAudioBuffer appends audio data to the input buffer
func (c *websocketsChat) AppendAudioBuffer(audioData []byte) error {
	// Encode audio data to base64
	encoded := base64.StdEncoding.EncodeToString(audioData)

	event := WebSocketInputAudioBufferAppendEvent{
		EventType: EventTypeInputAudioBufferAppend,
		Data: &WebSocketInputAudioBufferAppendEventData{
			Delta: encoded,
		},
	}

	return c.ws.sendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *websocketsChat) CompleteAudioBuffer() error {
	event := WebSocketInputAudioBufferCompleteEvent{
		EventType: EventTypeInputAudioBufferComplete,
	}

	return c.ws.sendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *websocketsChat) ClearAudioBuffer() error {
	event := WebSocketInputAudioBufferClearEvent{
		EventType: EventTypeInputAudioBufferClear,
	}

	return c.ws.sendEvent(event)
}

// CreateMessage creates a conversation message
func (c *websocketsChat) CreateMessage(content string) error {
	event := WebSocketConversationMessageCreateEvent{
		EventType: EventTypeConversationMessageCreate,
		Data: &WebSocketConversationMessageCreateEventData{
			Content: content,
		},
	}

	return c.ws.sendEvent(event)
}

// ClearConversation clears the conversation context
func (c *websocketsChat) ClearConversation() error {
	return c.ws.sendEvent(WebSocketConversationClearedEvent{
		baseWebSocketEvent: baseWebSocketEvent{
			EventType: EventTypeConversationClear,
		},
	})
}

// SubmitToolOutputs submits tool outputs for a chat
func (c *websocketsChat) SubmitToolOutputs(chatID string, toolOutputs []ToolOutput) error {
	event := WebSocketConversationChatSubmitToolOutputsEvent{
		EventType: EventTypeConversationChatSubmitToolOutputs,
		Data: &WebSocketConversationChatSubmitToolOutputsEventData{
			ChatID:      chatID,
			ToolOutputs: toolOutputs,
		},
	}

	return c.ws.sendEvent(event)
}

// CancelChat cancels the current chat
func (c *websocketsChat) CancelChat(chatID string) error {
	event := WebSocketConversationChatCancelEvent{
		baseWebSocketEvent: baseWebSocketEvent{
			EventType: EventTypeConversationChatCancel,
		},
		Data: &ConversationChatCancelData{
			ChatID: chatID,
		},
	}

	return c.ws.sendEvent(event)
}

// OnEvents registers multi events handler
func (c *websocketsChat) OnEvents(handlers map[WebSocketEventType]EventHandler) {
	for eventType, handler := range handlers {
		c.ws.OnEvent(eventType, handler)
	}
}

// OnEvent registers an event handler
func (c *websocketsChat) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

// WaitForChatCompleted waits for chat to complete
func (c *websocketsChat) WaitForChatCompleted(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		EventTypeConversationChatCompleted,
		EventTypeConversationChatFailed,
	}, timeout)
}

// WaitForChatCreated waits for chat to be created
func (c *websocketsChat) WaitForChatCreated(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		EventTypeConversationChatCreated,
	}, timeout)
}

// ChatEventHandler provides default handlers for chat events
type ChatEventHandler struct {
	OnChatCreated                          func(IWebSocketEvent) error
	OnChatUpdated                          func(IWebSocketEvent) error
	OnConversationChatCreated              func(*WebSocketConversationChatCreatedEvent) error
	OnConversationChatInProgress           func(IWebSocketEvent) error
	OnConversationMessageDelta             func(*WebSocketConversationMessageDeltaEvent) error
	OnConversationAudioSentenceStart       func(IWebSocketEvent) error
	OnConversationAudioDelta               func(*WebSocketConversationAudioDeltaEvent) error
	OnConversationMessageCompleted         func(IWebSocketEvent) error
	OnConversationAudioCompleted           func(IWebSocketEvent) error
	OnConversationChatCompleted            func(*WebSocketConversationChatCompletedEvent) error
	OnConversationChatFailed               func(IWebSocketEvent) error
	OnInputAudioBufferCompleted            func(IWebSocketEvent) error
	OnInputAudioBufferCleared              func(IWebSocketEvent) error
	OnConversationCleared                  func(IWebSocketEvent) error
	OnConversationChatCanceled             func(*WebSocketConversationChatCanceledEvent) error
	OnConversationAudioTranscriptUpdate    func(IWebSocketEvent) error
	OnConversationAudioTranscriptCompleted func(IWebSocketEvent) error
	OnConversationChatRequiresAction       func(*WebSocketConversationChatRequiresActionEvent) error
	OnInputAudioBufferSpeechStarted        func(IWebSocketEvent) error
	OnInputAudioBufferSpeechStopped        func(IWebSocketEvent) error
	OnError                                func(error) error
	OnClosed                               func() error
}

// RegisterHandlers registers all handlers with the client
func (h *ChatEventHandler) RegisterHandlers(client *websocketsChat) {
	if h.OnChatCreated != nil {
		client.OnEvent(EventTypeChatCreated, h.OnChatCreated)
	}

	if h.OnChatUpdated != nil {
		client.OnEvent(EventTypeChatUpdated, h.OnChatUpdated)
	}

	if h.OnConversationChatCreated != nil {
		client.OnEvent(EventTypeConversationChatCreated, func(event IWebSocketEvent) error {
			var chatEvent WebSocketConversationChatCreatedEvent
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
		client.OnEvent(EventTypeConversationMessageDelta, func(event IWebSocketEvent) error {
			var deltaEvent WebSocketConversationMessageDeltaEvent
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
		client.OnEvent(EventTypeConversationAudioDelta, func(event IWebSocketEvent) error {
			var audioEvent WebSocketConversationAudioDeltaEvent
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
		client.OnEvent(EventTypeConversationChatCompleted, func(event IWebSocketEvent) error {
			var completedEvent WebSocketConversationChatCompletedEvent
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
		client.OnEvent(EventTypeConversationChatCanceled, func(event IWebSocketEvent) error {
			var canceledEvent WebSocketConversationChatCanceledEvent
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
		client.OnEvent(EventTypeConversationChatRequiresAction, func(event IWebSocketEvent) error {
			var actionEvent WebSocketConversationChatRequiresActionEvent
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

type websocketChatBuilder struct {
	core *core
}

func newWebsocketChat(core *core) *websocketChatBuilder {
	return &websocketChatBuilder{
		core: core,
	}
}
