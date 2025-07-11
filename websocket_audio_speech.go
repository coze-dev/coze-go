package coze

import "context"

func (r *websocketAudioSpeechBuild) Create(ctx context.Context, req *CreateWebsocketAudioSpeechReq) *WebSocketAudioSpeech {
	return newWebSocketAudioSpeechClient(ctx, r.core, req)
}

type CreateWebsocketAudioSpeechReq struct {
	WebSocketClientOption *WebSocketClientOption
}

func (r *CreateWebsocketAudioSpeechReq) toQuery() map[string]string {
	q := map[string]string{}
	return q
}

type websocketAudioSpeechBuild struct {
	core *core
}

func newWebsocketAudioSpeechBuild(core *core) *websocketAudioSpeechBuild {
	return &websocketAudioSpeechBuild{
		core: core,
	}
}
