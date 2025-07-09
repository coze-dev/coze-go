package websockets

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

// WebSocketClient is the base WebSocket client
type WebSocketClient struct {
	baseURL     string
	path        string
	auth        Auth
	conn        *websocket.Conn
	sendChan    chan []byte
	receiveChan chan *WebSocketEvent
	closeChan   chan struct{}
	handlers    map[WebSocketEventType]EventHandler
	mu          sync.RWMutex
	connected   bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// Auth interface for authentication
type Auth interface {
	GetAuthHeader() (string, error)
}

// EventHandler represents a WebSocket event handler
type EventHandler func(event *WebSocketEvent) error

// WebSocketClientOption configures the WebSocket client
type WebSocketClientOption func(*WebSocketClient)

// WithEventHandler adds an event handler
func WithEventHandler(eventType WebSocketEventType, handler EventHandler) WebSocketClientOption {
	return func(c *WebSocketClient) {
		c.handlers[eventType] = handler
	}
}

// WithEventHandlers adds multiple event handlers
func WithEventHandlers(handlers map[WebSocketEventType]EventHandler) WebSocketClientOption {
	return func(c *WebSocketClient) {
		for eventType, handler := range handlers {
			c.handlers[eventType] = handler
		}
	}
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(baseURL, path string, auth Auth, opts ...WebSocketClientOption) *WebSocketClient {
	ctx, cancel := context.WithCancel(context.Background())

	client := &WebSocketClient{
		baseURL:     baseURL,
		path:        path,
		auth:        auth,
		sendChan:    make(chan []byte, 100),
		receiveChan: make(chan *WebSocketEvent, 100),
		closeChan:   make(chan struct{}),
		handlers:    make(map[WebSocketEventType]EventHandler),
		ctx:         ctx,
		cancel:      cancel,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Connect establishes the WebSocket connection
func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	// Build WebSocket URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	// Convert HTTP URL to WebSocket URL
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	u.Path = c.path

	// Get auth header
	authHeader, err := c.auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("failed to get auth header: %w", err)
	}

	// Setup headers
	headers := http.Header{}
	headers.Set("Authorization", authHeader)
	headers.Set("User-Agent", "coze-go/1.0")

	// Establish connection
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
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
func (c *WebSocketClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	c.cancel()

	// Close connection
	if c.conn != nil {
		c.conn.Close()
	}

	// Close channels
	close(c.closeChan)

	return nil
}

// IsConnected returns whether the client is connected
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SendEvent sends an event to the WebSocket
func (c *WebSocketClient) SendEvent(event interface{}) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
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
func (c *WebSocketClient) OnEvent(eventType WebSocketEventType, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[eventType] = handler
}

// WaitForEvent waits for specific events
func (c *WebSocketClient) WaitForEvent(eventTypes []WebSocketEventType, timeout time.Duration) (*WebSocketEvent, error) {
	eventChan := make(chan *WebSocketEvent, 1)

	// Register temporary handlers
	handlers := make(map[WebSocketEventType]EventHandler)
	for _, eventType := range eventTypes {
		handlers[eventType] = func(event *WebSocketEvent) error {
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
func (c *WebSocketClient) sendLoop() {
	for {
		select {
		case data := <-c.sendChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				// Handle error
				c.handleError(fmt.Errorf("failed to send message: %w", err))
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// receiveLoop handles receiving messages
func (c *WebSocketClient) receiveLoop() {
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

			var event WebSocketEvent
			if err := json.Unmarshal(message, &event); err != nil {
				c.handleError(fmt.Errorf("failed to unmarshal event: %w", err))
				continue
			}

			select {
			case c.receiveChan <- &event:
			default:
				// Channel full, skip event
			}
		}
	}
}

// handleEvents processes received events
func (c *WebSocketClient) handleEvents() {
	for {
		select {
		case event := <-c.receiveChan:
			c.mu.RLock()
			handler, exists := c.handlers[event.EventType]
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
func (c *WebSocketClient) handleError(err error) {
	c.mu.RLock()
	handler, exists := c.handlers[EventTypeError]
	c.mu.RUnlock()

	if exists && handler != nil {
		errorEvent := &WebSocketEvent{
			EventType: EventTypeError,
			Data:      []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		}
		handler(errorEvent)
	}
}
