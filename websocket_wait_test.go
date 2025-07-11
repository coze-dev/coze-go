package coze

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventWaiter_WaitOne(t *testing.T) {
	as := assert.New(t)

	t.Run("wait for one event", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)

		go func() {
			time.Sleep(100 * time.Millisecond)
			as.Nil(waiter.trigger(WebSocketEventTypeError))
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError}, false)
		as.Nil(err)
	})

	t.Run("wait for one event with timeout", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError}, false)
		as.NotNil(err)
	})
}

func TestEventWaiter_WaitAll(t *testing.T) {
	as := assert.New(t)
	t.Run("wait for all events", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)
		go func() {
			time.Sleep(50 * time.Millisecond)
			as.Nil(waiter.trigger(WebSocketEventTypeError))
			time.Sleep(50 * time.Millisecond)
			as.Nil(waiter.trigger(WebSocketEventTypeClientError))
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError, WebSocketEventTypeClientError}, true)
		as.Nil(err)
	})

	t.Run("wait for all events with timeout", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)
		go func() {
			time.Sleep(50 * time.Millisecond)
			as.Nil(waiter.trigger(WebSocketEventTypeError))
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError, WebSocketEventTypeClientError}, true)
		as.NotNil(err)
	})
}

func TestEventWaiter_WaitAny(t *testing.T) {
	as := assert.New(t)
	t.Run("wait for any event", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)
		go func() {
			time.Sleep(100 * time.Millisecond)
			as.Nil(waiter.trigger(WebSocketEventTypeError))
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError, WebSocketEventTypeClientError}, false)
		as.Nil(err)
	})

	t.Run("wait for any event with timeout", func(t *testing.T) {
		waiter := newEventWaiter(audioSpeechResponseEventTypes)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError, WebSocketEventTypeClientError}, false)
		as.NotNil(err)
	})
}

func TestEventWaiter_Shutdown(t *testing.T) {
	as := assert.New(t)
	waiter := newEventWaiter(audioSpeechResponseEventTypes)

	go func() {
		time.Sleep(100 * time.Millisecond)
		waiter.shutdown()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := waiter.wait(ctx, []WebSocketEventType{WebSocketEventTypeError, WebSocketEventTypeClientError}, true)
	as.Nil(err)
}
