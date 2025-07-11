package coze

import (
	"context"
	"testing"
)

func newMockCore() *core {
	return newCore(&clientOption{
		baseURL: "ws://localhost",
		auth:    &mockAuth{token: "mock-token"},
	})
}

func TestWebSocketChat_Basic(t *testing.T) {
	ctx := context.Background()
	core := newMockCore()
	req := &CreateWebsocketChatReq{}
	cli := newWebsocketChatClient(ctx, core, req)

	// 只测试不 panic，方法可调用
	_ = cli.Connect()
	_ = cli.IsConnected()
	_ = cli.Close()
}
