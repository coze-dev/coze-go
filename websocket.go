package coze

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
