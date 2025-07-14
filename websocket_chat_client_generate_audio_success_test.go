package coze

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/coze-dev/coze-go/examples/websockets/util"
	"github.com/stretchr/testify/assert"
)

type chatGenerateAudioSuccessTestdataHandler struct {
	BaseWebSocketChatHandler
	audio []byte
	mu    sync.Mutex
}

func (r *chatGenerateAudioSuccessTestdataHandler) OnConversationAudioDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioDeltaEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.audio = append(r.audio, event.Data.Content...)
	return nil
}

func (r *chatGenerateAudioSuccessTestdataHandler) assert(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	as := assert.New(t)

	f, err := os.CreateTemp("", "*-coze-ws-chat.wav")
	as.Nil(err)
	defer f.Close()
	defer os.Remove(f.Name())

	as.Nil(util.WritePCMToWavFile(f.Name(), r.audio))

	// audio
	expected, err := os.ReadFile("testdata/websocket_chat_generate_audio_success.wav")
	as.Nil(err)
	actual, err := os.ReadFile(f.Name())
	as.Nil(err)
	as.Equal(expected, actual)
}

func TestWebSocketChatGenerateAudioSuccess(t *testing.T) {
	as := assert.New(t)

	client := newWebsocketChatClient(context.Background(), newCore(&clientOption{
		baseURL:  CnBaseURL,
		logLevel: LogLevelDebug,
		logger:   newStdLogger(),
		auth:     NewTokenAuth("token"),
	}), &CreateWebsocketChatReq{
		WebSocketClientOption: &WebSocketClientOption{
			dial: connMockWebSocket(websocketChatGenerateAudioSuccessTestData),
		},
	})
	as.NotNil(client)
	handler := &chatGenerateAudioSuccessTestdataHandler{}
	client.RegisterHandler(handler)

	as.Nil(client.Connect())

	as.Nil(client.InputTextGenerateAudio(&WebSocketInputTextGenerateAudioEventData{
		Mode: WebSocketInputTextGenerateAudioModeText,
		Text: "亲，你怎么不说话了。",
	}))

	as.Nil(client.Wait(
		WebSocketEventTypeConversationAudioCompleted,
		WebSocketEventTypeError,
	))

	handler.assert(t)
}
