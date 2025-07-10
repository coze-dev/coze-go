package coze

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func (c *websocketChatBuilder) Create(ctx context.Context, req *CreateWebsocketChatReq) *websocketsChat {
	return newWebsocketChatClient(ctx, c.core, req)
}

type CreateWebsocketChatReq struct {
	// BotID is the ID of the bot.
	BotID string `json:"bot_id"`
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
		EventType: WebSocketEventTypeChatUpdate,
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
		EventType: WebSocketEventTypeInputAudioBufferAppend,
		Data: &WebSocketInputAudioBufferAppendEventData{
			Delta: encoded,
		},
	}

	return c.ws.sendEvent(event)
}

// CompleteAudioBuffer completes the audio buffer input
func (c *websocketsChat) CompleteAudioBuffer() error {
	event := WebSocketInputAudioBufferCompleteEvent{
		EventType: WebSocketEventTypeInputAudioBufferComplete,
	}

	return c.ws.sendEvent(event)
}

// ClearAudioBuffer clears the audio buffer
func (c *websocketsChat) ClearAudioBuffer() error {
	event := WebSocketInputAudioBufferClearEvent{
		EventType: WebSocketEventTypeInputAudioBufferClear,
	}

	return c.ws.sendEvent(event)
}

// CreateMessage creates a conversation message
func (c *websocketsChat) CreateMessage(content string) error {
	event := WebSocketConversationMessageCreateEvent{
		EventType: WebSocketEventTypeConversationMessageCreate,
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
			EventType: WebSocketEventTypeConversationClear,
		},
	})
}

// SubmitToolOutputs submits tool outputs for a chat
func (c *websocketsChat) SubmitToolOutputs(chatID string, toolOutputs []ToolOutput) error {
	event := WebSocketConversationChatSubmitToolOutputsEvent{
		EventType: WebSocketEventTypeConversationChatSubmitToolOutputs,
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
			EventType: WebSocketEventTypeConversationChatCancel,
		},
		Data: &ConversationChatCancelData{
			ChatID: chatID,
		},
	}

	return c.ws.sendEvent(event)
}

// WaitForChatCompleted waits for chat to complete
func (c *websocketsChat) WaitForChatCompleted(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeConversationChatCompleted,
		WebSocketEventTypeConversationChatFailed,
	}, timeout)
}

// WaitForChatCreated waits for chat to be created
func (c *websocketsChat) WaitForChatCreated(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeConversationChatCreated,
	}, timeout)
}

// OnEvent registers an event handler
func (c *websocketsChat) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

// OnEvent registers an event handler
func (c *websocketsChat) OnEvents(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

// OnChatCreated

func (c *websocketsChat) OnChatCreated(handler func(event *WebSocketChatCreatedEvent) error) {
	c.ws.OnEvent(WebSocketEventTypeChatCreated, func(event IWebSocketEvent) error {
		return handler(event.(*WebSocketChatCreatedEvent))
	})

}

// RegisterHandlers registers all handlers with the client
func (c *websocketsChat) RegisterHandlers(h IWebSocketChatHandler) {
	handlerType := reflect.TypeOf(h).Elem()
	handlerValue := reflect.ValueOf(h).Elem()

	for i := 0; i < handlerType.NumField(); i++ {
		field := handlerType.Field(i)
		methodValue := handlerValue.Field(i)
		if methodValue.IsNil() {
			continue
		}

		info, ok := websocketEvents[field.Name]
		if !ok {
			continue
		}
		// handler:=info.()

		c.OnEvent(info.get, func(event IWebSocketEvent) error {
			eventValue := reflect.ValueOf(event)
			if !eventValue.Type().AssignableTo(info.eventTypePtr) {
				return fmt.Errorf("invalid event type for %s", info.eventType)
			}
			result := methodValue.Call([]reflect.Value{eventValue})
			if len(result) > 0 && !result[0].IsNil() {
				return result[0].Interface().(error)
			}
			return nil
		})
	}

	if h.OnError != nil {
		c.OnEvent(WebSocketEventTypeError, func(ctx context.Context, cli *websocketsChat, event IWebSocketEvent) error {
			return h.OnError(ctx, event.(*WebSocketErrorEvent))
		})
	}
	if h.OnClientError != nil {
		c.OnEvent(WebSocketEventTypeClientError, func(event IWebSocketEvent) error {
			return h.OnClientError(event.(*WebSocketClientErrorEvent))
		})
	}
	if h.OnClosed != nil {
		c.OnEvent(WebSocketEventTypeClosed, func(event IWebSocketEvent) error {
			return h.OnClosed(event.(*WebSocketClosedEvent))
		})
	}
	if h.OnChatCreated != nil {
		c.OnEvent(WebSocketEventTypeChatCreated, func(event IWebSocketEvent) error {
			return h.OnChatCreated(event.(*WebSocketChatCreatedEvent))
		})
	}
	if h.OnChatUpdated != nil {
		c.OnEvent(WebSocketEventTypeChatUpdated, h.OnChatUpdated)
	}

	if h.OnConversationChatCreated != nil {
		c.OnEvent(WebSocketEventTypeConversationChatCreated, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeConversationChatInProgress, h.OnConversationChatInProgress)
	}

	if h.OnConversationMessageDelta != nil {
		c.OnEvent(WebSocketEventTypeConversationMessageDelta, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeConversationAudioSentenceStart, h.OnConversationAudioSentenceStart)
	}

	if h.OnConversationAudioDelta != nil {
		c.OnEvent(WebSocketEventTypeConversationAudioDelta, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeConversationMessageCompleted, h.OnConversationMessageCompleted)
	}

	if h.OnConversationAudioCompleted != nil {
		c.OnEvent(WebSocketEventTypeConversationAudioCompleted, h.OnConversationAudioCompleted)
	}

	if h.OnConversationChatCompleted != nil {
		c.OnEvent(WebSocketEventTypeConversationChatCompleted, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeConversationChatFailed, h.OnConversationChatFailed)
	}

	if h.OnInputAudioBufferCompleted != nil {
		c.OnEvent(WebSocketEventTypeInputAudioBufferCompleted, h.OnInputAudioBufferCompleted)
	}

	if h.OnInputAudioBufferCleared != nil {
		c.OnEvent(WebSocketEventTypeInputAudioBufferCleared, h.OnInputAudioBufferCleared)
	}

	if h.OnConversationCleared != nil {
		c.OnEvent(WebSocketEventTypeConversationCleared, h.OnConversationCleared)
	}

	if h.OnConversationChatCanceled != nil {
		c.OnEvent(WebSocketEventTypeConversationChatCanceled, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeConversationAudioTranscriptUpdate, h.OnConversationAudioTranscriptUpdate)
	}

	if h.OnConversationAudioTranscriptCompleted != nil {
		c.OnEvent(WebSocketEventTypeConversationAudioTranscriptCompleted, h.OnConversationAudioTranscriptCompleted)
	}

	if h.OnConversationChatRequiresAction != nil {
		c.OnEvent(WebSocketEventTypeConversationChatRequiresAction, func(event IWebSocketEvent) error {
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
		c.OnEvent(WebSocketEventTypeInputAudioBufferSpeechStarted, h.OnInputAudioBufferSpeechStarted)
	}

	if h.OnInputAudioBufferSpeechStopped != nil {
		c.OnEvent(WebSocketEventTypeInputAudioBufferSpeechStopped, h.OnInputAudioBufferSpeechStopped)
	}

	if h.OnError != nil {
		c.OnEvent(WebSocketEventTypeError, func(event IWebSocketEvent) error {
			return h.OnError(fmt.Errorf("WebSocket error: %s", string(event.Data)))
		})
	}

	if h.OnClosed != nil {
		c.OnEvent(WebSocketEventTypeClosed, func(event IWebSocketEvent) error {
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
