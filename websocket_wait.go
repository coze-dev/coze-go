package coze

import (
	"context"
	"sync"
)

type eventWaiter struct {
	events sync.Map
}

type eventWaitState struct {
	triggered bool
	once      sync.Once
	ch        chan struct{}
}

func newEventWaiter() *eventWaiter {
	return &eventWaiter{
		events: sync.Map{},
	}
}

func (r *eventWaiter) getState(key string) *eventWaitState {
	state, _ := r.events.LoadOrStore(key, &eventWaitState{
		ch: make(chan struct{}, 1),
	})
	return state.(*eventWaitState)
}

func (r *eventWaiter) shutdown() {
	r.events.Range(func(key, value any) bool {
		state := value.(*eventWaitState)
		state.once.Do(func() {
			close(state.ch)
		})
		return true
	})
}

// waitAll: true 表示等待所有事件都触发，false 表示等待任意一个事件触发
func (r *eventWaiter) wait(ctx context.Context, keys []string, waitAll bool) error {
	if len(keys) == 0 {
		return nil
	}

	// 如果只有一个事件，直接调用单个等待方法
	if len(keys) == 1 {
		return r.waitOne(ctx, keys[0])
	}

	// 获取所有事件的状态
	states := make([]*eventWaitState, len(keys))
	for i, key := range keys {
		states[i] = r.getState(key)
	}

	if waitAll {
		return r.waitAll(ctx, states)
	} else {
		return r.waitAny(ctx, states)
	}
}

// wait 等待单个事件
func (r *eventWaiter) waitOne(ctx context.Context, key string) error {
	state := r.getState(key)

	if state.triggered {
		return nil
	}

	select {
	case <-state.ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// waitAll 等待所有事件都触发
func (r *eventWaiter) waitAll(ctx context.Context, states []*eventWaitState) error {
	// 检查是否所有事件都已经触发
	allTriggered := true
	for _, state := range states {
		if !state.triggered {
			allTriggered = false
			break
		}
	}
	if allTriggered {
		return nil
	}

	// 创建一个用于聚合的 channel
	done := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	completed := 0

	// 为每个未触发的事件启动一个 goroutine
	for _, state := range states {
		if state.triggered {
			completed++
			continue
		}

		wg.Add(1)
		go func(s *eventWaitState) {
			defer wg.Done()
			<-s.ch

			mu.Lock()
			completed++
			if completed == len(states) {
				close(done)
			}
			mu.Unlock()
		}(state)
	}

	// 如果所有事件在启动 goroutine 后都已完成，直接关闭 done
	mu.Lock()
	if completed == len(states) {
		close(done)
	}
	mu.Unlock()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// waitAny 等待任意一个事件触发
func (r *eventWaiter) waitAny(ctx context.Context, states []*eventWaitState) error {
	// 检查是否已有事件触发
	for _, state := range states {
		if state.triggered {
			return nil
		}
	}

	// 创建一个用于聚合的 channel
	done := make(chan struct{}, 1)

	// 为每个事件启动一个 goroutine
	for _, state := range states {
		go func(s *eventWaitState) {
			<-s.ch
			select {
			case done <- struct{}{}:
			default:
			}
		}(state)
	}

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *eventWaiter) trigger(key string) {
	state := r.getState(key)

	state.once.Do(func() {
		state.triggered = true
		close(state.ch) // 关闭 channel 来通知所有等待者
	})
}
