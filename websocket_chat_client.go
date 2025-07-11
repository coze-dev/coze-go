package coze

import (
	"context"
	"time"
)

type WebSocketChat struct {
	ctx  context.Context
	core *core
	req  *CreateWebsocketChatReq

	ws *websocketClient
}

func newWebsocketChatClient(ctx context.Context, core *core, req *CreateWebsocketChatReq) *WebSocketChat {
	ws := newWebSocketClient(
		&WebSocketClientOption{
			ctx:   ctx,
			core:  core,
			path:  "/v1/chat",
			query: req.toQuery(),
		},
	)

	client := &WebSocketChat{
		ctx:  ctx,
		core: core,
		req:  req,
		ws:   ws,
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *WebSocketChat) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *WebSocketChat) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *WebSocketChat) IsConnected() bool {
	return c.ws.IsConnected()
}

func (c *WebSocketChat) ChatUpdate(data *WebSocketChatUpdateEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeChatUpdate, data).(WebSocketChatUpdateEvent))
}
func (c *WebSocketChat) InputAudioBufferAppend(data *WebSocketInputAudioBufferAppendEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferAppend, data).(WebSocketInputAudioBufferAppendEvent))
}
func (c *WebSocketChat) InputAudioBufferComplete(data *WebSocketInputAudioBufferCompleteEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferComplete, data).(WebSocketInputAudioBufferCompleteEvent))
}
func (c *WebSocketChat) InputAudioBufferClear(data *WebSocketInputAudioBufferClearEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferClear, data).(WebSocketInputAudioBufferClearEvent))
}
func (c *WebSocketChat) ConversationMessageCreate(data *WebSocketConversationMessageCreateEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeConversationMessageCreate, data).(WebSocketConversationMessageCreateEvent))
}
func (c *WebSocketChat) ConversationClear(data *WebSocketConversationClearEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeConversationClear, data).(WebSocketConversationClearEvent))
}
func (c *WebSocketChat) ConversationChatSubmitToolOutputs(data *WebSocketConversationChatSubmitToolOutputsEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeConversationChatSubmitToolOutputs, data).(WebSocketConversationChatSubmitToolOutputsEvent))
}
func (c *WebSocketChat) ConversationChatCancel(data *WebSocketConversationChatCancelEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeConversationChatCancel, data).(WebSocketConversationChatCancelEvent))
}

// Wait waits for chat to complete
func (c *WebSocketChat) Wait(timeout time.Duration) (IWebSocketEvent, error) {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeConversationChatCompleted,
		WebSocketEventTypeConversationChatFailed,
	}, timeout)
}

// OnEvent registers an event handler
func (c *WebSocketChat) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

func registerEventHandler[T any](c *WebSocketChat, eventType WebSocketEventType, handler func(ctx context.Context, cli *WebSocketChat, event *T) error) {
	c.ws.OnEvent(eventType, func(event IWebSocketEvent) error {
		return handler(c.ctx, c, (any)(event).(*T))
	})
}

func (c *WebSocketChat) OnError(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketErrorEvent) error) {
	registerEventHandler(c, WebSocketEventTypeError, handler)
}

func (c *WebSocketChat) OnClientError(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketClientErrorEvent) error) {
	registerEventHandler(c, WebSocketEventTypeClientError, handler)
}

func (c *WebSocketChat) OnClosed(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketClosedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeClosed, handler)
}

func (c *WebSocketChat) OnChatCreated(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketChatCreatedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeChatCreated, handler)
}

func (c *WebSocketChat) OnChatUpdated(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketChatUpdatedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeChatUpdated, handler)
}

func (c *WebSocketChat) OnConversationChatCreated(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCreatedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatCreated, handler)
}

func (c *WebSocketChat) OnConversationChatInProgress(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatInProgressEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatInProgress, handler)
}

func (c *WebSocketChat) OnConversationMessageDelta(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageDeltaEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationMessageDelta, handler)
}

func (c *WebSocketChat) OnConversationAudioSentenceStart(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioSentenceStartEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationAudioSentenceStart, handler)
}

func (c *WebSocketChat) OnConversationAudioDelta(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioDeltaEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationAudioDelta, handler)
}

func (c *WebSocketChat) OnConversationMessageCompleted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageCompletedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationMessageCompleted, handler)
}

func (c *WebSocketChat) OnConversationAudioCompleted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioCompletedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationAudioCompleted, handler)
}

func (c *WebSocketChat) OnConversationChatCompleted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCompletedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatCompleted, handler)
}

func (c *WebSocketChat) OnConversationChatFailed(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatFailedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatFailed, handler)
}

func (c *WebSocketChat) OnInputAudioBufferCompleted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferCompletedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeInputAudioBufferCompleted, handler)
}

func (c *WebSocketChat) OnInputAudioBufferCleared(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferClearedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeInputAudioBufferCleared, handler)
}

func (c *WebSocketChat) OnConversationCleared(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationClearedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationCleared, handler)
}

func (c *WebSocketChat) OnConversationChatCanceled(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCanceledEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatCanceled, handler)
}

func (c *WebSocketChat) OnConversationAudioTranscriptUpdate(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptUpdateEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationAudioTranscriptUpdate, handler)
}

func (c *WebSocketChat) OnConversationAudioTranscriptCompleted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptCompletedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationAudioTranscriptCompleted, handler)
}

func (c *WebSocketChat) OnConversationChatRequiresAction(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatRequiresActionEvent) error) {
	registerEventHandler(c, WebSocketEventTypeConversationChatRequiresAction, handler)
}

func (c *WebSocketChat) OnInputAudioBufferSpeechStarted(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStartedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeInputAudioBufferSpeechStarted, handler)
}

func (c *WebSocketChat) OnInputAudioBufferSpeechStopped(handler func(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStoppedEvent) error) {
	registerEventHandler(c, WebSocketEventTypeInputAudioBufferSpeechStopped, handler)
}

// RegisterHandler registers all handlers with the client
func (c *WebSocketChat) RegisterHandler(h IWebSocketChatHandler) {
	c.OnClientError(h.OnClientError)
	c.OnClosed(h.OnClosed)
	c.OnError(h.OnError)
	c.OnChatCreated(h.OnChatCreated)
	c.OnChatUpdated(h.OnChatUpdated)
	c.OnConversationChatCreated(h.OnConversationChatCreated)
	c.OnConversationChatInProgress(h.OnConversationChatInProgress)
	c.OnConversationMessageDelta(h.OnConversationMessageDelta)
	c.OnConversationAudioSentenceStart(h.OnConversationAudioSentenceStart)
	c.OnConversationAudioDelta(h.OnConversationAudioDelta)
	c.OnConversationMessageCompleted(h.OnConversationMessageCompleted)
	c.OnConversationAudioCompleted(h.OnConversationAudioCompleted)
	c.OnConversationChatCompleted(h.OnConversationChatCompleted)
	c.OnConversationChatFailed(h.OnConversationChatFailed)
	c.OnInputAudioBufferCompleted(h.OnInputAudioBufferCompleted)
	c.OnInputAudioBufferCleared(h.OnInputAudioBufferCleared)
	c.OnConversationCleared(h.OnConversationCleared)
	c.OnConversationChatCanceled(h.OnConversationChatCanceled)
	c.OnConversationAudioTranscriptUpdate(h.OnConversationAudioTranscriptUpdate)
	c.OnConversationAudioTranscriptCompleted(h.OnConversationAudioTranscriptCompleted)
	c.OnConversationChatRequiresAction(h.OnConversationChatRequiresAction)
	c.OnInputAudioBufferSpeechStarted(h.OnInputAudioBufferSpeechStarted)
	c.OnInputAudioBufferSpeechStopped(h.OnInputAudioBufferSpeechStopped)
}
