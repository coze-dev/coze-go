package coze

type websocketAudio struct {
	core *core
}

func newWebsocketAudio(core *core) *websocketAudio {
	return &websocketAudio{
		core: core,
	}
}
