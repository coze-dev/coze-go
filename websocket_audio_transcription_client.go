package coze

import (
	"context"
)

var _ WebSocketClient = (*WebSocketAudioTranscription)(nil)

type WebSocketAudioTranscription struct {
	ctx  context.Context
	core *core

	ws *websocketClient
}

func newWebSocketAudioTranscriptionClient(ctx context.Context, core *core, req *CreateWebsocketAudioTranscriptionReq) *WebSocketAudioTranscription {
	ws := newWebSocketClient(mergeWebSocketClientOption(req.WebSocketClientOption, &WebSocketClientOption{
		ctx:                ctx,
		core:               core,
		path:               "/v1/audio/transcriptions",
		query:              req.toQuery(),
		responseEventTypes: audioTranscriptionResponseEventTypes,
	}))

	return &WebSocketAudioTranscription{
		ctx:  ctx,
		core: core,
		ws:   ws,
	}
}

// Connect establishes the WebSocket connection
func (c *WebSocketAudioTranscription) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *WebSocketAudioTranscription) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *WebSocketAudioTranscription) IsConnected() bool {
	return c.ws.IsConnected()
}

func (c *WebSocketAudioTranscription) TranscriptionsUpdate(data *WebSocketTranscriptionsUpdateEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeTranscriptionsUpdate, data).(WebSocketTranscriptionsUpdateEvent))
}

func (c *WebSocketAudioTranscription) InputAudioBufferAppend(data *WebSocketInputAudioBufferAppendEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferAppend, data).(WebSocketInputAudioBufferAppendEvent))
}

func (c *WebSocketAudioTranscription) InputAudioBufferComplete(data *WebSocketInputAudioBufferCompleteEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferComplete, data).(WebSocketInputAudioBufferCompleteEvent))
}

func (c *WebSocketAudioTranscription) InputAudioBufferClear(data *WebSocketInputAudioBufferClearEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputAudioBufferClear, data).(WebSocketInputAudioBufferClearEvent))
}

// Wait waits for transcription to complete
func (c *WebSocketAudioTranscription) Wait() error {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeTranscriptionsMessageCompleted,
		WebSocketEventTypeError,
	}, false)
}

// OnEvent registers an event handler
func (c *WebSocketAudioTranscription) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

func registerAudioTranscriptionEventHandler[T any](c *WebSocketAudioTranscription, eventType WebSocketEventType, handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *T) error) {
	c.ws.OnEvent(eventType, func(event IWebSocketEvent) error {
		return handler(c.ctx, c, (any)(event).(*T))
	})
}

func (c *WebSocketAudioTranscription) OnClientError(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeClientError, handler)
}

func (c *WebSocketAudioTranscription) OnError(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeError, handler)
}

func (c *WebSocketAudioTranscription) OnClosed(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeClosed, handler)
}

func (c *WebSocketAudioTranscription) OnTranscriptionsCreated(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeTranscriptionsCreated, handler)
}

func (c *WebSocketAudioTranscription) OnTranscriptionsUpdated(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeTranscriptionsUpdated, handler)
}

func (c *WebSocketAudioTranscription) OnInputAudioBufferCompleted(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeInputAudioBufferCompleted, handler)
}

func (c *WebSocketAudioTranscription) OnInputAudioBufferCleared(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeInputAudioBufferCleared, handler)
}

func (c *WebSocketAudioTranscription) OnTranscriptionsMessageUpdate(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeTranscriptionsMessageUpdate, handler)
}

func (c *WebSocketAudioTranscription) OnTranscriptionsMessageCompleted(handler func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error) {
	registerAudioTranscriptionEventHandler(c, WebSocketEventTypeTranscriptionsMessageCompleted, handler)
}

func (c *WebSocketAudioTranscription) RegisterHandler(h IWebSocketAudioTranscriptionHandler) {
	c.OnClientError(h.OnClientError)
	c.OnClosed(h.OnClosed)
	c.OnError(h.OnError)
	c.OnTranscriptionsCreated(h.OnTranscriptionsCreated)
	c.OnTranscriptionsUpdated(h.OnTranscriptionsUpdated)
	c.OnInputAudioBufferCompleted(h.OnInputAudioBufferCompleted)
	c.OnInputAudioBufferCleared(h.OnInputAudioBufferCleared)
	c.OnTranscriptionsMessageUpdate(h.OnTranscriptionsMessageUpdate)
	c.OnTranscriptionsMessageCompleted(h.OnTranscriptionsMessageCompleted)
}
