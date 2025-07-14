package coze

import (
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed testdata/websocket_speech_success.txt
var websocketSpeechSuccessTestData string

//go:embed testdata/websocket_transcriptions_success.txt
var websocketTranscriptionsSuccessTestData string

//go:embed testdata/websocket_chat_success.txt
var websocketChatSuccessTestData string

//go:embed testdata/websocket_chat_generate_audio_success.txt
var websocketChatGenerateAudioSuccessTestData string

type testdataWebSocketItem struct {
	Type  string // send, receive
	Event string //
}

func readTestdataWebSocket(content string) ([]*testdataWebSocketItem, error) {
	l := strings.Split(content, "\n")
	res := []*testdataWebSocketItem{}

	typ := ""
	for _, v := range l {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if typ == "" {
			typ = v
			continue
		}
		if typ != "send" && typ != "receive" {
			return nil, errors.New("type not send or receive")
		}
		res = append(res, &testdataWebSocketItem{
			Type:  typ,
			Event: v,
		})
		typ = ""
	}
	if typ != "" {
		return nil, errors.New("type not empty")
	}
	return res, nil
}

func connMockWebSocket(mockData string) func(dialer websocket.Dialer, urlStr string, requestHeader http.Header) (websocketConn, error) {
	items, err := readTestdataWebSocket(mockData)
	if err != nil {
		panic(err)
	}
	return func(dialer websocket.Dialer, urlStr string, requestHeader http.Header) (websocketConn, error) {
		return newMockWebSocketConn(items), nil
	}
}

type mockWebSocketConn struct {
	mu        sync.RWMutex
	items     []*testdataWebSocketItem
	idx       int
	sendCh    chan string
	receiveCh chan string
}

func newMockWebSocketConn(expectedItems []*testdataWebSocketItem) *mockWebSocketConn {
	conn := &mockWebSocketConn{
		items:     expectedItems,
		sendCh:    make(chan string),
		receiveCh: make(chan string),
	}
	go func() {
		for _, v := range expectedItems {
			if v.Type == "send" {
				conn.sendCh <- v.Event
			} else if v.Type == "receive" {
				conn.receiveCh <- v.Event
			} else {
				panic("invalid type " + v.Type)
			}
		}
	}()
	return conn
}

func (r *mockWebSocketConn) Close() error {
	return nil
}

func (r *mockWebSocketConn) readCh(ch chan string) (string, error) {
	for {
		excepted := ""
		read := false
		select {
		case excepted = <-ch:
			read = true
		default:
		}

		if read {
			r.mu.Lock()
			r.idx++
			r.mu.Unlock()
			return excepted, nil
		}

		r.mu.RLock()
		left := r.idx < len(r.items)
		r.mu.RUnlock()
		if !left {
			return "", fmt.Errorf("no left")
		}
		time.Sleep(time.Millisecond)
	}
}

func (r *mockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	excepted, err := r.readCh(r.sendCh)
	if err != nil {
		return fmt.Errorf("no left, write failed: %s", data)
	}

	if excepted != string(data) {
		return fmt.Errorf("excepted data not match, excepted: %q, actual: %q", excepted, string(data))
	}
	return nil
}

func (r *mockWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	excepted, err := r.readCh(r.receiveCh)
	if err != nil {
		return 0, nil, net.ErrClosed
	}

	return websocket.TextMessage, []byte(excepted), nil
}
