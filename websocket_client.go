package coze

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// websocketClient is the base WebSocket client
type websocketClient struct {
	opt *WebSocketClientOption

	core        *core
	conn        *websocket.Conn
	sendChan    chan []byte          // 发送队列, 长度 100
	receiveChan chan IWebSocketEvent // 接收队列, 长度 100
	closeChan   chan struct{}
	processing  sync.WaitGroup
	handlers    sync.Map // map[WebSocketEventType]EventHandler
	mu          sync.RWMutex
	connected   bool
	ctx         context.Context
	cancel      context.CancelFunc
	waiter      *eventWaiter
}

type WebSocketClientOption struct {
	ctx                 context.Context
	core                *core
	path                string
	query               map[string]string
	SendChanCapacity    int           // 默认 1000
	ReceiveChanCapacity int           // 默认 1000
	HandshakeTimeout    time.Duration // 默认 3s
}

// EventHandler represents a WebSocket event handler
type EventHandler func(event IWebSocketEvent) error

// newWebSocketClient creates a new WebSocket client
func newWebSocketClient(opt *WebSocketClientOption) *websocketClient {
	ctx, cancel := context.WithCancel(context.Background())

	if opt.ReceiveChanCapacity == 0 {
		opt.ReceiveChanCapacity = 1000
	}
	if opt.SendChanCapacity == 0 {
		opt.SendChanCapacity = 1000
	}
	if opt.HandshakeTimeout == 0 {
		opt.HandshakeTimeout = 3 * time.Second
	}

	client := &websocketClient{
		opt:         opt,
		core:        opt.core,
		sendChan:    make(chan []byte, opt.SendChanCapacity),
		receiveChan: make(chan IWebSocketEvent, opt.ReceiveChanCapacity),
		closeChan:   make(chan struct{}),
		handlers:    sync.Map{},
		ctx:         ctx,
		cancel:      cancel,
		waiter:      newEventWaiter(),
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *websocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	baseURL := c.opt.core.baseURL
	path := c.opt.path
	auth := c.opt.core.auth
	query := c.opt.query

	// Build WebSocket URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	// Convert HTTP URL to WebSocket URL
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}
	if u.Host == "api.coze.cn" {
		u.Host = "ws.coze.cn"
	} else if u.Host == "api.coze.com" {
		u.Host = "ws.coze.com"
	}

	u.Path = path

	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	// Get auth header
	accessToken, err := auth.Token(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth header: %w", err)
	}

	// Setup headers
	headers := http.Header{}
	// auth
	headers.Set("Authorization", "Bearer "+accessToken)
	// agent
	headers.Set("User-Agent", userAgent)
	headers.Set("X-Coze-Client-User-Agent", clientUserAgent)

	// Establish connection
	dialer := websocket.Dialer{
		HandshakeTimeout: c.opt.HandshakeTimeout,
	}

	c.core.Log(c.ctx, LogLevelDebug, "[%s] connecting to websocket: %s", c.opt.path, u.String())
	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	c.connected = true

	// Start goroutines
	go c.sendLoop()
	go c.receiveLoop()
	go c.handleEvents()

	return nil
}

// Close closes the WebSocket connection
func (c *websocketClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	// wait for receive channels to be empty
	c.processing.Wait()

	c.connected = false
	c.cancel()

	// Close connection
	var err error
	if c.conn != nil {
		err = c.conn.Close()
	}

	// Close channels
	close(c.closeChan)

	return err
}

// IsConnected returns whether the client is connected
func (c *websocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// 发送事件
func (c *websocketClient) sendEvent(event any) error {
	if !c.IsConnected() {
		return fmt.Errorf("websocket not connected")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	select {
	case c.sendChan <- data:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("context cancelled")
	default:
		return fmt.Errorf("send channel full")
	}
}

// OnEvent registers an event handler
func (c *websocketClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.handlers.Store(eventType, handler)
}

// WaitForEvent waits for specific events
func (c *websocketClient) WaitForEvent(eventTypes []WebSocketEventType, waitAll bool) error {
	keys := make([]string, 0, 10)
	for _, eventType := range eventTypes {
		keys = append(keys, string(eventType))
	}
	return c.waiter.wait(c.ctx, keys, waitAll)
}

// sendLoop handles sending messages
func (c *websocketClient) sendLoop() {
	for {
		select {
		case data := <-c.sendChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.handleError(WebSocketEventTypeClientError, fmt.Errorf("failed to send message: %w", err))
				continue
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// receiveLoop handles receiving messages
func (c *websocketClient) receiveLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				c.handleError(WebSocketEventTypeClientError, fmt.Errorf("failed to read message: %w", err))
				c.waiter.shutdown()
				return
			}

			event, err := parseWebSocketEvent(message)
			if err != nil {
				// TODO: 这里不应该发 client error？
				c.handleError(WebSocketEventTypeClientError, err)
				continue
			}

			c.waiter.trigger(string(event.GetEventType()))

			if event.GetEventType() == WebSocketEventTypeSpeechAudioUpdate {
				c.core.Log(c.ctx, LogLevelDebug, "[%s] receive event, type=%s, event=%s", c.opt.path, event.GetEventType(), event.(*WebSocketSpeechAudioUpdateEvent).dumpWithoutBinary())
			} else if event.GetEventType() == WebSocketEventTypeConversationAudioDelta {
				c.core.Log(c.ctx, LogLevelDebug, "[%s] receive event, type=%s, event=%s", c.opt.path, event.GetEventType(), event.(*WebSocketConversationAudioDeltaEvent).dumpWithoutBinary())
			} else {
				c.core.Log(c.ctx, LogLevelDebug, "[%s] receive event, type=%s, event=%s", c.opt.path, event.GetEventType(), message)
			}

			// 没有 timeout 或者 channel full 处理, 暂时符合预期
			c.processing.Add(1)
			c.receiveChan <- event
		}
	}
}

// handleEvents processes received events
func (c *websocketClient) handleEvents() {
	for {
		select {
		case event := <-c.receiveChan:
			c.handleEvent(event)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *websocketClient) handleEvent(event IWebSocketEvent) {
	defer c.processing.Done()

	handler := c.getHandler(event.GetEventType())

	if handler != nil {
		if err := handler(event); err != nil {
			// TODO: handler 返回错误类型？
			c.handleError(WebSocketEventTypeClientError, fmt.Errorf("event handler error: %w", err))
		}
	}
}

// handleError handles errors
func (c *websocketClient) handleError(eventType WebSocketEventType, err error) {
	c.core.Log(c.ctx, LogLevelWarn, "[%s] receive event, type=%s, event=%s", c.opt.path, eventType, err)

	handler := c.getHandler(eventType)
	if handler == nil {
		return
	}
	if eventType == WebSocketEventTypeClientError {
		handler(&WebSocketClientErrorEvent{
			baseWebSocketEvent: baseWebSocketEvent{
				EventType: eventType,
			},
			Data: err,
		})
	}
}

func (c *websocketClient) getHandler(eventType WebSocketEventType) EventHandler {
	handler, ok := c.handlers.Load(eventType)
	if !ok {
		return nil
	}
	return handler.(EventHandler)
}
