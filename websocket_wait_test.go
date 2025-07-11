package coze

import (
	"context"
	"testing"
	"time"
)

func TestEventWaiter_WaitOne(t *testing.T) {
	waiter := newEventWaiter()

	t.Run("wait for one event", func(t *testing.T) {
		go func() {
			time.Sleep(100 * time.Millisecond)
			waiter.trigger("event1")
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event1"}, false)
		if err != nil {
			t.Errorf("wait() error = %v, wantErr nil", err)
		}
	})

	t.Run("wait for one event with timeout", func(t *testing.T) {
		waiter := newEventWaiter()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event2"}, false)
		if err == nil {
			t.Errorf("wait() error = %v, wantErr context.DeadlineExceeded", err)
		}
	})
}

func TestEventWaiter_WaitAll(t *testing.T) {
	waiter := newEventWaiter()

	t.Run("wait for all events", func(t *testing.T) {
		go func() {
			time.Sleep(50 * time.Millisecond)
			waiter.trigger("event1")
			time.Sleep(50 * time.Millisecond)
			waiter.trigger("event2")
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event1", "event2"}, true)
		if err != nil {
			t.Errorf("wait() error = %v, wantErr nil", err)
		}
	})

	t.Run("wait for all events with timeout", func(t *testing.T) {
		waiter := newEventWaiter()
		go func() {
			time.Sleep(50 * time.Millisecond)
			waiter.trigger("event1")
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event1", "event2"}, true)
		if err == nil {
			t.Errorf("wait() error = %v, wantErr context.DeadlineExceeded", err)
		}
	})
}

func TestEventWaiter_WaitAny(t *testing.T) {
	waiter := newEventWaiter()

	t.Run("wait for any event", func(t *testing.T) {
		go func() {
			time.Sleep(100 * time.Millisecond)
			waiter.trigger("event1")
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event1", "event2"}, false)
		if err != nil {
			t.Errorf("wait() error = %v, wantErr nil", err)
		}
	})

	t.Run("wait for any event with timeout", func(t *testing.T) {
		waiter := newEventWaiter()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waiter.wait(ctx, []string{"event1", "event2"}, false)
		if err == nil {
			t.Errorf("wait() error = %v, wantErr context.DeadlineExceeded", err)
		}
	})
}

func TestEventWaiter_Shutdown(t *testing.T) {
	waiter := newEventWaiter()

	go func() {
		time.Sleep(100 * time.Millisecond)
		waiter.shutdown()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := waiter.wait(ctx, []string{"event1", "event2"}, true)
	if err != nil {
		t.Errorf("wait() error = %v, wantErr nil", err)
	}
}