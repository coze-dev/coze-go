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
	dial        websocketDialer
	conn        websocketConn
	sendChan    chan IWebSocketEvent // 发送队列, 长度 1000
	receiveChan chan IWebSocketEvent // 接收队列, 长度 1000
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
	responseEventTypes  []WebSocketEventType
	dial                websocketDialer
	SendChanCapacity    int           // 默认 1000
	ReceiveChanCapacity int           // 默认 1000
	HandshakeTimeout    time.Duration // 默认 3s
}

type websocketConn interface {
	Close() error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

type websocketDialer func(dialer websocket.Dialer, urlStr string, requestHeader http.Header) (websocketConn, error)

func dialWebSocket(dialer websocket.Dialer, urlStr string, requestHeader http.Header) (websocketConn, error) {
	conn, _, err := dialer.Dial(urlStr, requestHeader)
	return conn, err
}

func mergeWebSocketClientOption(opt *WebSocketClientOption, other *WebSocketClientOption) *WebSocketClientOption {
	if opt == nil {
		opt = &WebSocketClientOption{}
	}
	opt.ctx = other.ctx
	opt.core = other.core
	opt.path = other.path
	opt.query = other.query
	opt.responseEventTypes = other.responseEventTypes
	return opt
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
	if opt.dial == nil {
		opt.dial = dialWebSocket
	}

	client := &websocketClient{
		opt:         opt,
		core:        opt.core,
		dial:        opt.dial,
		sendChan:    make(chan IWebSocketEvent, opt.SendChanCapacity),
		receiveChan: make(chan IWebSocketEvent, opt.ReceiveChanCapacity),
		closeChan:   make(chan struct{}),
		handlers:    sync.Map{},
		ctx:         ctx,
		cancel:      cancel,
		waiter:      newEventWaiter(opt.responseEventTypes),
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
	conn, err := c.dial(dialer, u.String(), headers)
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
func (c *websocketClient) sendEvent(event IWebSocketEvent) error {
	if !c.IsConnected() {
		return fmt.Errorf("websocket not connected")
	}

	select {
	case c.sendChan <- event:
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
	return c.waiter.wait(c.ctx, eventTypes, waitAll)
}

// sendLoop handles sending messages
func (c *websocketClient) sendLoop() {
	for {
		select {
		case event := <-c.sendChan:
			data, err := json.Marshal(event)
			if err != nil {
				if c.core.logLevel <= LogLevelDebug {
					c.core.Log(c.ctx, LogLevelDebug, "[%s] send event, marshal_failed, type=%s, event=%s, err=%s", c.opt.path, event.GetEventType(), mustToJson(event), err)
				}
				c.handleClientError(fmt.Errorf("failed to marshal event: %w", err))
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				if c.core.logLevel <= LogLevelDebug {
					c.core.Log(c.ctx, LogLevelDebug, "[%s] send event, write_failed, type=%s, event=%s, err=%s", c.opt.path, event.GetEventType(), mustToJson(event), err)
				}
				c.handleClientError(fmt.Errorf("failed to send message: %w", err))
				continue
			}
			if c.core.logLevel <= LogLevelDebug {
				c.core.Log(c.ctx, LogLevelDebug, "[%s] send event, type=%s, event=%s", c.opt.path, event.GetEventType(), mustToJson(event))
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
				c.handleClientError(fmt.Errorf("failed to read message: %w", err))
				c.waiter.shutdown()
				return
			}

			event, err := parseWebSocketEvent(message)
			if err != nil {
				c.core.Log(c.ctx, LogLevelDebug, "[%s] receive event, parse_failed, event=%s, err=%s", c.opt.path, message, err)
				c.handleClientError(err)
				continue
			}

			if err := c.waiter.trigger(event.GetEventType()); err != nil {
				c.core.Log(c.ctx, LogLevelWarn, "[%s] trigger event failed, event_type=%s, err=%s", c.opt.path, event.GetEventType(), err)
			}

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
			c.core.Log(c.ctx, LogLevelWarn, "[%s] handler %s failed, logid=%s, err=%s", c.opt.path, event.GetEventType(), event.GetDetail().LogID, err)
		}
	}
}

// handleClientError handles errors
func (c *websocketClient) handleClientError(err error) {
	handler := c.getHandler(WebSocketEventTypeClientError)
	if handler == nil || err == nil {
		return
	}
	if err := handler(&WebSocketClientErrorEvent{
		baseWebSocketEvent: baseWebSocketEvent{
			EventType: WebSocketEventTypeClientError,
		},
		Data: err,
	}); err != nil {
		c.core.Log(c.ctx, LogLevelWarn, "[%s] handler %s failed, err=%s", c.opt.path, WebSocketEventTypeClientError, err)
	}
}

func (c *websocketClient) getHandler(eventType WebSocketEventType) EventHandler {
	handler, ok := c.handlers.Load(eventType)
	if !ok {
		return nil
	}
	return handler.(EventHandler)
}
