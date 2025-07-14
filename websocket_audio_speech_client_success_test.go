package coze

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/coze-dev/coze-go/examples/websockets/util"
	"github.com/stretchr/testify/assert"
)

type speechSuccessTestdataHandler struct {
	BaseWebSocketAudioSpeechHandler
	audio []byte
	mu    sync.Mutex
}

func (r *speechSuccessTestdataHandler) OnSpeechAudioUpdate(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.audio = append(r.audio, event.Data.Delta...)
	return nil
}

func (r *speechSuccessTestdataHandler) assert(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	as := assert.New(t)

	f, err := os.CreateTemp("", "*-coze-ws-speech.wav")
	as.Nil(err)
	defer f.Close()
	defer os.Remove(f.Name())

	as.Nil(util.WritePCMToWavFile(f.Name(), r.audio))

	expected, err := os.ReadFile("testdata/websocket_speech_success.wav")
	as.Nil(err)
	actual, err := os.ReadFile(f.Name())
	as.Nil(err)
	as.Equal(expected, actual)
}

func TestWebSocketSpeechSuccess(t *testing.T) {
	as := assert.New(t)

	client := newWebSocketAudioSpeechClient(context.Background(), newCore(&clientOption{
		baseURL:  CnBaseURL,
		logLevel: LogLevelDebug,
		logger:   newStdLogger(),
		auth:     NewTokenAuth("token"),
	}), &CreateWebsocketAudioSpeechReq{
		WebSocketClientOption: &WebSocketClientOption{
			dial: connMockWebSocket(websocketSpeechSuccessTestData),
		},
	})
	as.NotNil(client)
	handler := &speechSuccessTestdataHandler{}
	client.RegisterHandler(handler)

	as.Nil(client.Connect())

	as.Nil(client.InputTextBufferAppend(&WebSocketInputTextBufferAppendEventData{
		Delta: "今天天气不错",
	}))
	as.Nil(client.InputTextBufferComplete(nil))
	as.Nil(client.Wait())

	handler.assert(t)
}
