package coze

import (
	_ "embed"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

//go:embed testdata/websocket_speech_success.txt
var websocketSpeechSuccessTestData string

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
