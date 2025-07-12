package coze

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func mockDialWebSocket(dialer websocket.Dialer, urlStr string, requestHeader http.Header) (websocketConn, error) {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &mockWebSocketConn{
		path: urlParsed.Path,
	}, nil
}

type mockWebSocketConn struct {
	path string
}

func (r *mockWebSocketConn) Close() error {
	// TODO implement me
	panic("implement me")
}

func (r *mockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	// TODO implement me
	panic("implement me")
}

func (r *mockWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	// TODO implement me
	panic("implement me")
}
