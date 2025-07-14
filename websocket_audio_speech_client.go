package coze

import (
	"context"
)

var _ WebSocketClient = (*WebSocketAudioSpeech)(nil)

type WebSocketAudioSpeech struct {
	ctx  context.Context
	core *core
	ws   *websocketClient
}

func newWebSocketAudioSpeechClient(ctx context.Context, core *core, req *CreateWebsocketAudioSpeechReq) *WebSocketAudioSpeech {
	ws := newWebSocketClient(mergeWebSocketClientOption(req.WebSocketClientOption, &WebSocketClientOption{
		ctx:                ctx,
		core:               core,
		path:               "/v1/audio/speech",
		query:              req.toQuery(),
		responseEventTypes: audioSpeechResponseEventTypes,
	}))

	return &WebSocketAudioSpeech{
		ctx:  ctx,
		core: core,
		ws:   ws,
	}
}

// Connect establishes the WebSocket connection
func (c *WebSocketAudioSpeech) Connect() error {
	return c.ws.Connect()
}

// Close closes the WebSocket connection
func (c *WebSocketAudioSpeech) Close() error {
	return c.ws.Close()
}

// IsConnected returns whether the client is connected
func (c *WebSocketAudioSpeech) IsConnected() bool {
	return c.ws.IsConnected()
}

func (c *WebSocketAudioSpeech) SpeechUpdate(data *WebSocketSpeechUpdateEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeSpeechUpdate, data).(WebSocketSpeechUpdateEvent))
}

// InputTextBufferAppend appends text to the input text buffer
func (c *WebSocketAudioSpeech) InputTextBufferAppend(data *WebSocketInputTextBufferAppendEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputTextBufferAppend, data).(WebSocketInputTextBufferAppendEvent))
}

// InputTextBufferComplete completes the input text buffer
func (c *WebSocketAudioSpeech) InputTextBufferComplete(data *WebSocketInputTextBufferCompleteEventData) error {
	return c.ws.sendEvent(newWebSocketEvent(WebSocketEventTypeInputTextBufferComplete, data).(WebSocketInputTextBufferCompleteEvent))
}

// Wait waits for speech audio to complete
func (c *WebSocketAudioSpeech) Wait() error {
	return c.ws.WaitForEvent([]WebSocketEventType{
		WebSocketEventTypeSpeechAudioCompleted,
		WebSocketEventTypeError,
	}, false)
}

// OnEvent registers an event handler
func (c *WebSocketAudioSpeech) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.ws.OnEvent(eventType, handler)
}

func registerAudioSpeechEventHandler[T any](c *WebSocketAudioSpeech, eventType WebSocketEventType, handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *T) error) {
	c.ws.OnEvent(eventType, func(event IWebSocketEvent) error {
		return handler(c.ctx, c, (any)(event).(*T))
	})
}

func (c *WebSocketAudioSpeech) OnClientError(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeClientError, handler)
}

func (c *WebSocketAudioSpeech) OnError(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeError, handler)
}

func (c *WebSocketAudioSpeech) OnClosed(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeClosed, handler)
}

func (c *WebSocketAudioSpeech) OnSpeechCreated(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeSpeechCreated, handler)
}

func (c *WebSocketAudioSpeech) OnSpeechUpdated(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeSpeechUpdated, handler)
}

func (c *WebSocketAudioSpeech) OnInputTextBufferCompleted(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeInputTextBufferCompleted, handler)
}

func (c *WebSocketAudioSpeech) OnSpeechAudioUpdate(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeSpeechAudioUpdate, handler)
}

func (c *WebSocketAudioSpeech) OnSpeechAudioCompleted(handler func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error) {
	registerAudioSpeechEventHandler(c, WebSocketEventTypeSpeechAudioCompleted, handler)
}

func (c *WebSocketAudioSpeech) RegisterHandler(h IWebSocketAudioSpeechHandler) {
	c.OnClientError(h.OnClientError)
	c.OnClosed(h.OnClosed)
	c.OnError(h.OnError)
	c.OnSpeechCreated(h.OnSpeechCreated)
	c.OnSpeechUpdated(h.OnSpeechUpdated)
	c.OnInputTextBufferCompleted(h.OnInputTextBufferCompleted)
	c.OnSpeechAudioUpdate(h.OnSpeechAudioUpdate)
	c.OnSpeechAudioCompleted(h.OnSpeechAudioCompleted)
}
