package coze

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWebSocketAudioSpeechHandler struct {
	OnClientErrorFunc              func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error
	OnClosedFunc                   func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error
	OnErrorFunc                    func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error
	OnSpeechCreatedFunc            func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error
	OnSpeechUpdatedFunc            func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error
	OnInputTextBufferCompletedFunc func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error
	OnSpeechAudioUpdateFunc        func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error
	OnSpeechAudioCompletedFunc     func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error
}

func (h *mockWebSocketAudioSpeechHandler) OnClientError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error {
	if h.OnClientErrorFunc != nil {
		return h.OnClientErrorFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnClosed(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error {
	if h.OnClosedFunc != nil {
		return h.OnClosedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error {
	if h.OnErrorFunc != nil {
		return h.OnErrorFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnSpeechCreated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error {
	if h.OnSpeechCreatedFunc != nil {
		return h.OnSpeechCreatedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnSpeechUpdated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error {
	if h.OnSpeechUpdatedFunc != nil {
		return h.OnSpeechUpdatedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnInputTextBufferCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error {
	if h.OnInputTextBufferCompletedFunc != nil {
		return h.OnInputTextBufferCompletedFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnSpeechAudioUpdate(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error {
	if h.OnSpeechAudioUpdateFunc != nil {
		return h.OnSpeechAudioUpdateFunc(ctx, cli, event)
	}
	return nil
}

func (h *mockWebSocketAudioSpeechHandler) OnSpeechAudioCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error {
	if h.OnSpeechAudioCompletedFunc != nil {
		return h.OnSpeechAudioCompletedFunc(ctx, cli, event)
	}
	return nil
}

func TestWebSocketAudioSpeech(t *testing.T) {
	t.Run("IsConnected-false", func(t *testing.T) {
		client := newWebSocketAudioSpeechClient(context.Background(), &core{}, &CreateWebsocketAudioSpeechReq{})
		assert.NotNil(t, client)
		assert.False(t, client.IsConnected())
	})

	t.Run("SendEvent-NotConnected", func(t *testing.T) {
		client := newWebSocketAudioSpeechClient(context.Background(), &core{}, &CreateWebsocketAudioSpeechReq{})
		assert.NotNil(t, client)

		err := client.SpeechUpdate(&WebSocketSpeechUpdateEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())

		err = client.InputTextBufferAppend(&WebSocketInputTextBufferAppendEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())

		err = client.InputTextBufferComplete(&WebSocketInputTextBufferCompleteEventData{})
		assert.Error(t, err)
		assert.Equal(t, "websocket not connected", err.Error())
	})

	t.Run("RegisterHandler", func(t *testing.T) {
		client := newWebSocketAudioSpeechClient(context.Background(), &core{}, &CreateWebsocketAudioSpeechReq{})
		assert.NotNil(t, client)

		var (
			onClientErrorCalled              bool
			onClosedCalled                   bool
			onErrorCalled                    bool
			onSpeechCreatedCalled            bool
			onSpeechUpdatedCalled            bool
			onInputTextBufferCompletedCalled bool
			onSpeechAudioUpdateCalled        bool
			onSpeechAudioCompletedCalled     bool
		)

		handler := &mockWebSocketAudioSpeechHandler{
			OnClientErrorFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error {
				onClientErrorCalled = true
				return nil
			},
			OnClosedFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error {
				onClosedCalled = true
				return nil
			},
			OnErrorFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error {
				onErrorCalled = true
				return nil
			},
			OnSpeechCreatedFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error {
				onSpeechCreatedCalled = true
				return nil
			},
			OnSpeechUpdatedFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error {
				onSpeechUpdatedCalled = true
				return nil
			},
			OnInputTextBufferCompletedFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error {
				onInputTextBufferCompletedCalled = true
				return nil
			},
			OnSpeechAudioUpdateFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error {
				onSpeechAudioUpdateCalled = true
				return nil
			},
			OnSpeechAudioCompletedFunc: func(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error {
				onSpeechAudioCompletedCalled = true
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

		triggerHandler(WebSocketEventTypeSpeechCreated, &WebSocketSpeechCreatedEvent{})
		assert.True(t, onSpeechCreatedCalled)

		triggerHandler(WebSocketEventTypeSpeechUpdated, &WebSocketSpeechUpdatedEvent{})
		assert.True(t, onSpeechUpdatedCalled)

		triggerHandler(WebSocketEventTypeInputTextBufferCompleted, &WebSocketInputTextBufferCompletedEvent{})
		assert.True(t, onInputTextBufferCompletedCalled)

		triggerHandler(WebSocketEventTypeSpeechAudioUpdate, &WebSocketSpeechAudioUpdateEvent{})
		assert.True(t, onSpeechAudioUpdateCalled)

		triggerHandler(WebSocketEventTypeSpeechAudioCompleted, &WebSocketSpeechAudioCompletedEvent{})
		assert.True(t, onSpeechAudioCompletedCalled)
	})
}
