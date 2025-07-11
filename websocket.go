package coze

type WebSocketClient interface {
	Connect() error
	Close() error
	IsConnected() bool
	Wait() error
	OnEvent(eventType WebSocketEventType, handler EventHandler)
}

// websockets is the main WebSocket client that provides access to all WebSocket services
type websockets struct {
	core  *core
	Audio *websocketAudio
	Chat  *websocketChatBuilder
}

func newWebSockets(core *core) *websockets {
	return &websockets{
		core:  core,
		Audio: newWebsocketAudio(core),
		Chat:  newWebsocketChat(core),
	}
}
