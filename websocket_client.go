package coze

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// websocketClient is the base WebSocket client
type websocketClient struct {
	opt *WebSocketClientOption

	conn        *websocket.Conn
	sendChan    chan []byte          // 发送队列, 长度 100
	receiveChan chan IWebSocketEvent // 接收队列, 长度 100
	closeChan   chan struct{}
	handlers    map[WebSocketEventType]EventHandler
	mu          sync.RWMutex
	connected   bool
	ctx         context.Context
	cancel      context.CancelFunc
}

type WebSocketClientOption struct {
	ctx                 context.Context
	core                *core
	path                string
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
		sendChan:    make(chan []byte, opt.SendChanCapacity),
		receiveChan: make(chan IWebSocketEvent, opt.ReceiveChanCapacity),
		closeChan:   make(chan struct{}),
		handlers:    make(map[WebSocketEventType]EventHandler),
		ctx:         ctx,
		cancel:      cancel,
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

	u.Path = path

	// Get auth header
	accessToken, err := auth.Token(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth header: %w", err)
	}

	// Setup headers
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+accessToken)
	headers.Set("User-Agent", "coze-go/1.0") // todo

	// Establish connection
	dialer := websocket.Dialer{
		HandshakeTimeout: c.opt.HandshakeTimeout,
	}

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
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[eventType] = handler
}

// WaitForEvent waits for specific events
func (c *websocketClient) WaitForEvent(eventTypes []WebSocketEventType, timeout time.Duration) (IWebSocketEvent, error) {
	eventChan := make(chan IWebSocketEvent, 1)

	// Register temporary handlers
	handlers := make(map[WebSocketEventType]EventHandler)
	for _, eventType := range eventTypes {
		handlers[eventType] = func(event IWebSocketEvent) error {
			select {
			case eventChan <- event:
			default:
			}
			return nil
		}
	}

	c.mu.Lock()
	originalHandlers := make(map[WebSocketEventType]EventHandler)
	for eventType, handler := range handlers {
		originalHandlers[eventType] = c.handlers[eventType]
		c.handlers[eventType] = handler
	}
	c.mu.Unlock()

	// Wait for event or timeout
	select {
	case event := <-eventChan:
		// Restore original handlers
		c.mu.Lock()
		for eventType, handler := range originalHandlers {
			if handler != nil {
				c.handlers[eventType] = handler
			} else {
				delete(c.handlers, eventType)
			}
		}
		c.mu.Unlock()
		return event, nil
	case <-time.After(timeout):
		// Restore original handlers
		c.mu.Lock()
		for eventType, handler := range originalHandlers {
			if handler != nil {
				c.handlers[eventType] = handler
			} else {
				delete(c.handlers, eventType)
			}
		}
		c.mu.Unlock()
		return nil, fmt.Errorf("timeout waiting for event")
	case <-c.ctx.Done():
		return nil, fmt.Errorf("context cancelled")
	}
}

// sendLoop handles sending messages
func (c *websocketClient) sendLoop() {
	for {
		select {
		case data := <-c.sendChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.handleError(fmt.Errorf("failed to send message: %w", err))
				return
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
				c.handleError(fmt.Errorf("failed to read message: %w", err))
				return
			}

			event, err := parseWebSocketEvent(message)
			if err := json.Unmarshal(message, &event); err != nil {
				c.handleError(err)
				continue
			}

			select {
			case c.receiveChan <- event:
			default:
				// todo log
				// Channel full, skip event
			}
		}
	}
}

// handleEvents processes received events
func (c *websocketClient) handleEvents() {
	for {
		select {
		case event := <-c.receiveChan:
			c.mu.RLock()
			handler, exists := c.handlers[event.GetEventType()]
			c.mu.RUnlock()

			if exists && handler != nil {
				if err := handler(event); err != nil {
					c.handleError(fmt.Errorf("event handler error: %w", err))
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// handleError handles errors
func (c *websocketClient) handleError(err error) {
	c.mu.RLock()
	handler, ok := c.handlers[EventTypeError]
	c.mu.RUnlock()

	if ok && handler != nil {
		errorEvent := &WebSocketErrorEvent{
			baseWebSocketEvent: baseWebSocketEvent{
				EventType: EventTypeError,
			},
			// todo
			Data: &ErrorData{
				Code: 0,
				Msg:  "",
			},
		}
		handler(errorEvent)
	}
}
