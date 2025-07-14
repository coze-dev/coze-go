package coze

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/coze-dev/coze-go/examples/websockets/util"
	"github.com/stretchr/testify/assert"
)

type chatSuccessTestdataHandler struct {
	BaseWebSocketChatHandler
	audio []byte
	text  string
	mu    sync.Mutex
}

func (r *chatSuccessTestdataHandler) OnConversationMessageDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageDeltaEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.text += event.Data.Content
	return nil
}

func (r *chatSuccessTestdataHandler) OnConversationAudioDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioDeltaEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.audio = append(r.audio, event.Data.Content...)
	return nil
}

func (r *chatSuccessTestdataHandler) assert(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	as := assert.New(t)

	f, err := os.CreateTemp("", "*-coze-ws-chat.wav")
	as.Nil(err)
	defer f.Close()
	defer os.Remove(f.Name())

	as.Nil(util.WritePCMToWavFile(f.Name(), r.audio))

	// audio
	expected, err := os.ReadFile("testdata/websocket_chat_success.wav")
	as.Nil(err)
	fmt.Println(f.Name())
	actual, err := os.ReadFile(f.Name())
	as.Nil(err)
	as.Equal(expected, actual)

	// text
	as.Equal("是啊，好天气总能让人心情也跟着变好呢！你有没有打算趁着这好天气出门走走，做点有意思的事儿？  ", r.text)
}

func TestWebSocketChatSuccess(t *testing.T) {
	as := assert.New(t)

	client := newWebsocketChatClient(context.Background(), newCore(&clientOption{
		baseURL:  CnBaseURL,
		logLevel: LogLevelDebug,
		logger:   newStdLogger(),
		auth:     NewTokenAuth("token"),
	}), &CreateWebsocketChatReq{
		WebSocketClientOption: &WebSocketClientOption{
			dial: connMockWebSocket(websocketChatSuccessTestData),
		},
	})
	as.NotNil(client)
	handler := &chatSuccessTestdataHandler{}
	client.RegisterHandler(handler)

	as.Nil(client.Connect())

	as.Nil(client.ConversationMessageCreate(&WebSocketConversationMessageCreateEventData{
		Role:        MessageRoleUser,
		ContentType: MessageContentTypeText,
		Content:     "今天天气真不错",
	}))

	as.Nil(client.Wait())

	handler.assert(t)
}
