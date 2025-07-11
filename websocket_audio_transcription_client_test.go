package coze

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWebSocketAudioTranscriptionHandler struct {
	OnClientErrorFunc                    func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error
	OnClosedFunc                         func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error
	OnErrorFunc                          func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error
	OnTranscriptionsCreatedFunc          func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error
	OnTranscriptionsUpdatedFunc          func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error
	OnInputAudioBufferCompletedFunc      func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error
	OnInputAudioBufferClearedFunc        func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error
	OnTranscriptionsMessageUpdateFunc    func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error
	OnTranscriptionsMessageCompletedFunc func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error
}

func (h *mockWebSocketAudioTranscriptionHandler) OnClientError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error {
	if h.OnClientErrorFunc != nil {
		return h.OnClientErrorFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnClosed(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error {
	if h.OnClosedFunc != nil {
		return h.OnClosedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error {
	if h.OnErrorFunc != nil {
		return h.OnErrorFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnTranscriptionsCreated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error {
	if h.OnTranscriptionsCreatedFunc != nil {
		return h.OnTranscriptionsCreatedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnTranscriptionsUpdated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error {
	if h.OnTranscriptionsUpdatedFunc != nil {
		return h.OnTranscriptionsUpdatedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnInputAudioBufferCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error {
	if h.OnInputAudioBufferCompletedFunc != nil {
		return h.OnInputAudioBufferCompletedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnInputAudioBufferCleared(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error {
	if h.OnInputAudioBufferClearedFunc != nil {
		return h.OnInputAudioBufferClearedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnTranscriptionsMessageUpdate(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error {
	if h.OnTranscriptionsMessageUpdateFunc != nil {
		return h.OnTranscriptionsMessageUpdateFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioTranscriptionHandler) OnTranscriptionsMessageCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error {
	if h.OnTranscriptionsMessageCompletedFunc != nil {
		return h.OnTranscriptionsMessageCompletedFunc(ctx, cli, event)
	}
	return nil
}

func TestWebSocketAudioTranscription(t *testing.T) {
	t.Run("IsConnected-false", func(t *testing.T) {
		client := newWebSocketAudioTranscriptionClient(context.Background(), &core{}, &CreateWebsocketAudioTranscriptionReq{})
		assert.NotNil(t, client)
		assert.False(t, client.IsConnected())
	})

	t.Run("SendEvent-NotConnected", func(t *testing.T) {
		client := newWebSocketAudioTranscriptionClient(context.Background(), &core{}, &CreateWebsocketAudioTranscriptionReq{})
		assert.NotNil(t, client)

		err := client.TranscriptionsUpdate(&WebSocketTranscriptionsUpdateEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())

		err = client.InputAudioBufferAppend(&WebSocketInputAudioBufferAppendEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())

		err = client.InputAudioBufferComplete(&WebSocketInputAudioBufferCompleteEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())

		err = client.InputAudioBufferClear(&WebSocketInputAudioBufferClearEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())
	})

	t.Run("RegisterHandler", func(t *testing.T) {
		client := newWebSocketAudioTranscriptionClient(context.Background(), &core{}, &CreateWebsocketAudioTranscriptionReq{})
		assert.NotNil(t, client)

		var (
			onClientErrorCalled                    bool
			onClosedCalled                         bool
			onErrorCalled                          bool
			onTranscriptionsCreatedCalled          bool
			onTranscriptionsUpdatedCalled          bool
			onInputAudioBufferCompletedCalled      bool
			onInputAudioBufferClearedCalled        bool
			onTranscriptionsMessageUpdateCalled    bool
			onTranscriptionsMessageCompletedCalled bool
		)

		handler := &mockWebSocketAudioTranscriptionHandler{
			OnClientErrorFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error {
				onClientErrorCalled = true
				return nil
			},
			OnClosedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error {
				onClosedCalled = true
				return nil
			},
			OnErrorFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error {
				onErrorCalled = true
				return nil
			},
			OnTranscriptionsCreatedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error {
				onTranscriptionsCreatedCalled = true
				return nil
			},
			OnTranscriptionsUpdatedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error {
				onTranscriptionsUpdatedCalled = true
				return nil
			},
			OnInputAudioBufferCompletedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error {
				onInputAudioBufferCompletedCalled = true
				return nil
			},
			OnInputAudioBufferClearedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error {
				onInputAudioBufferClearedCalled = true
				return nil
			},
			OnTranscriptionsMessageUpdateFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error {
				onTranscriptionsMessageUpdateCalled = true
				return nil
			},
			OnTranscriptionsMessageCompletedFunc: func(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error {
				onTranscriptionsMessageCompletedCalled = true
				return nil
			},
		}
		client.RegisterHandler(handler)

		// Helper to trigger handlers
		triggerHandler := func(eventType WebSocketEventType, event IWebSocketEvent) {
			h, ok := client.ws.handlers.Load(eventType)
			assert.True(t, ok)
			err := h.(EventHandler)(event)
			assert.NoError(t, err)
		}

		triggerHandler(WebSocketEventTypeClientError, &WebSocketClientErrorEvent{})
		assert.True(t, onClientErrorCalled)

		triggerHandler(WebSocketEventTypeClosed, &WebSocketClosedEvent{})
		assert.True(t, onClosedCalled)

		triggerHandler(WebSocketEventTypeError, &WebSocketErrorEvent{})
		assert.True(t, onErrorCalled)

		triggerHandler(WebSocketEventTypeTranscriptionsCreated, &WebSocketTranscriptionsCreatedEvent{})
		assert.True(t, onTranscriptionsCreatedCalled)

		triggerHandler(WebSocketEventTypeTranscriptionsUpdated, &WebSocketTranscriptionsUpdatedEvent{})
		assert.True(t, onTranscriptionsUpdatedCalled)

		triggerHandler(WebSocketEventTypeInputAudioBufferCompleted, &WebSocketInputAudioBufferCompletedEvent{})
		assert.True(t, onInputAudioBufferCompletedCalled)

		triggerHandler(WebSocketEventTypeInputAudioBufferCleared, &WebSocketInputAudioBufferClearedEvent{})
		assert.True(t, onInputAudioBufferClearedCalled)

		triggerHandler(WebSocketEventTypeTranscriptionsMessageUpdate, &WebSocketTranscriptionsMessageUpdateEvent{})
		assert.True(t, onTranscriptionsMessageUpdateCalled)

		triggerHandler(WebSocketEventTypeTranscriptionsMessageCompleted, &WebSocketTranscriptionsMessageCompletedEvent{})
		assert.True(t, onTranscriptionsMessageCompletedCalled)
	})
}
