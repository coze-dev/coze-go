package coze

import "context"

func (r *websocketAudioTranscriptionBuild) Create(ctx context.Context, req *CreateWebsocketAudioTranscriptionReq) *WebSocketAudioTranscription {
	return newWebSocketAudioTranscriptionClient(ctx, r.core, req)
}

type CreateWebsocketAudioTranscriptionReq struct {
	WebSocketClientOption *WebSocketClientOption
}

func (r *CreateWebsocketAudioTranscriptionReq) toQuery() map[string]string {
	q := map[string]string{}
	return q
}

type websocketAudioTranscriptionBuild struct {
	core *core
}

func newWebsocketAudioTranscriptionBuild(core *core) *websocketAudioTranscriptionBuild {
	return &websocketAudioTranscriptionBuild{
		core: core,
	}
}
