package coze

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type transcriptionsSuccessTestdataHandler struct {
	BaseWebSocketAudioTranscriptionHandler
	content string
	mu      sync.Mutex
}

func (r *transcriptionsSuccessTestdataHandler) OnTranscriptionsMessageUpdate(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.content += event.Data.Content
	return nil
}

func (r *transcriptionsSuccessTestdataHandler) assert(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	as := assert.New(t)

	as.Equal("今天天气不错。", r.content)
}

func TestWebSocketTranscriptionsSuccess(t *testing.T) {
	as := assert.New(t)

	client := newWebSocketAudioTranscriptionClient(context.Background(), newCore(&clientOption{
		baseURL:  CnBaseURL,
		logLevel: LogLevelDebug,
		logger:   newStdLogger(),
		auth:     NewTokenAuth("token"),
	}), &CreateWebsocketAudioTranscriptionReq{
		WebSocketClientOption: &WebSocketClientOption{
			dial: connMockWebSocket(websocketTranscriptionsSuccessTestData),
		},
	})
	as.NotNil(client)
	handler := &transcriptionsSuccessTestdataHandler{}
	client.RegisterHandler(handler)

	as.Nil(client.Connect())

	audioData, err := os.ReadFile("testdata/websocket_speech_success.wav")
	as.Nil(err)

	as.Nil(client.InputAudioBufferAppend(&WebSocketInputAudioBufferAppendEventData{
		Delta: audioData,
	}))
	as.Nil(client.InputAudioBufferComplete(nil))
	as.Nil(client.Wait())

	handler.assert(t)
}
