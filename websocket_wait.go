package coze

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

const maxEventSize = websocketEventTypeSize

type eventWaiter struct {
	eventMask      uint64              // 位图
	eventTriggered [maxEventSize]int32 // 1: 已触发, 0: 未触发
	eventChannels  [maxEventSize]chan struct{}
	eventOnce      [maxEventSize]sync.Once
	enableEvents   []WebSocketEventType
}

func getWebSocketEventTypeIndex(WebSocketEventType WebSocketEventType) (int, bool) {
	i, ok := websocketEventTypeIndex[WebSocketEventType]
	return i, ok
}

func newEventWaiter(enableEvents []WebSocketEventType) *eventWaiter {
	oew := &eventWaiter{
		enableEvents: enableEvents,
	}
	for _, v := range enableEvents {
		i, ok := getWebSocketEventTypeIndex(v)
		if !ok {
			continue
		}
		oew.eventChannels[i] = make(chan struct{})
	}
	return oew
}

func (oew *eventWaiter) wait(ctx context.Context, eventTypes []WebSocketEventType, waitAll bool) error {
	if len(eventTypes) <= 0 {
		return nil
	} else if len(eventTypes) == 1 {
		return oew.waitOne(ctx, eventTypes[0])
	}
	if waitAll {
		return oew.waitAll(ctx, eventTypes)
	}
	_, err := oew.waitAny(ctx, eventTypes)
	return err
}

// waitOne 等待单个事件
func (oew *eventWaiter) waitOne(ctx context.Context, eventType WebSocketEventType) error {
	i, ok := getWebSocketEventTypeIndex(eventType)
	if !ok {
		return fmt.Errorf("wait event_type: %s not found", eventType)
	}
	// 检查事件是否已经触发
	if atomic.LoadInt32(&oew.eventTriggered[i]) == 1 {
		return nil
	}

	select {
	case <-oew.eventChannels[i]:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// waitAny 等待任意一个事件触发
func (oew *eventWaiter) waitAny(ctx context.Context, eventTypes []WebSocketEventType) (WebSocketEventType, error) {
	if len(eventTypes) == 0 {
		return "", errors.New("no event types specified")
	}

	// 创建事件类型掩码
	var eventMask uint64
	for _, eventType := range eventTypes {
		i, ok := getWebSocketEventTypeIndex(eventType)
		if !ok {
			return "", fmt.Errorf("wait event_type: %s not found", eventType)
		}
		eventMask |= 1 << i
	}

	// 检查是否已有事件触发
	currentMask := atomic.LoadUint64(&oew.eventMask)
	if triggeredMask := currentMask & eventMask; triggeredMask != 0 {
		// 找到第一个触发的事件
		for _, eventType := range eventTypes {
			i, ok := getWebSocketEventTypeIndex(eventType)
			if !ok {
				return "", fmt.Errorf("wait event_type: %s not found", eventType)
			}
			if triggeredMask&(1<<i) != 0 {
				return eventType, nil
			}
		}
	}

	// 创建一个聚合 channel
	done := make(chan WebSocketEventType, 1)
	var wg sync.WaitGroup

	for _, et := range eventTypes {
		wg.Add(1)
		go func(eventType WebSocketEventType) {
			defer wg.Done()

			// 再次检查是否已触发（避免竞态条件）
			i, _ := getWebSocketEventTypeIndex(eventType) // 不需要检查合法性
			if atomic.LoadInt32(&oew.eventTriggered[i]) == 1 {
				select {
				case done <- eventType:
				default:
				}
				return
			}

			select {
			case <-oew.eventChannels[i]:
				select {
				case done <- eventType:
				default:
				}
			case <-ctx.Done():
			}
		}(et)
	}

	// 等待第一个事件或超时
	select {
	case eventType := <-done:
		return eventType, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// waitAll 等待所有指定事件都触发
func (oew *eventWaiter) waitAll(ctx context.Context, WebSocketEventTypes []WebSocketEventType) error {
	if len(WebSocketEventTypes) == 0 {
		return nil
	}

	// 创建事件类型掩码
	var targetMask uint64
	for _, eventType := range WebSocketEventTypes {
		i, ok := getWebSocketEventTypeIndex(eventType)
		if !ok {
			return fmt.Errorf("wait event_type: %s not found", eventType)
		}
		targetMask |= 1 << i
	}

	// 检查是否所有事件都已触发
	checkAllTriggered := func() bool {
		currentMask := atomic.LoadUint64(&oew.eventMask)
		return (currentMask & targetMask) == targetMask
	}

	if checkAllTriggered() {
		return nil
	}

	// 等待所有事件
	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, et := range WebSocketEventTypes {
		wg.Add(1)
		go func(eventType WebSocketEventType) {
			defer wg.Done()

			i, _ := getWebSocketEventTypeIndex(eventType) // 不需要检查合法性
			if atomic.LoadInt32(&oew.eventTriggered[i]) == 0 {
				select {
				case <-oew.eventChannels[i]:
				case <-ctx.Done():
					return
				}
			}
		}(et)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if checkAllTriggered() {
			return nil
		}
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TriggerEvent 触发事件
func (oew *eventWaiter) trigger(eventType WebSocketEventType) error {
	i, ok := getWebSocketEventTypeIndex(eventType)
	if !ok {
		return fmt.Errorf("wait event_type: %s not found", eventType)
	}

	// 使用 sync.Once 确保只触发一次
	oew.eventOnce[i].Do(func() {
		atomic.StoreInt32(&oew.eventTriggered[i], 1)

		// 原子性地设置位掩码
		for {
			old := atomic.LoadUint64(&oew.eventMask)
			new := old | (1 << i)
			if atomic.CompareAndSwapUint64(&oew.eventMask, old, new) {
				break
			}
		}

		safeCloseChan(oew.eventChannels[i]) // 关闭 channel 通知所有等待者
	})

	return nil
}

func (oew *eventWaiter) shutdown() {
	// 关闭所有 channel
	for _, eventType := range oew.enableEvents {
		i, ok := getWebSocketEventTypeIndex(eventType)
		if !ok {
			continue
		}

		safeCloseChan(oew.eventChannels[i])
	}
}

func safeCloseChan[T any](ch chan T) {
	defer func() {
		_ = recover()
	}()
	close(ch)
}
