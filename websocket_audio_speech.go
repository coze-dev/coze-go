package coze

import "context"

func (r *websocketAudioSpeechBuild) Create(ctx context.Context, req *CreateWebsocketAudioSpeechReq) *WebSocketAudioSpeech {
	return newWebSocketAudioSpeechClient(ctx, r.core, req)
}

type CreateWebsocketAudioSpeechReq struct {
	// BotID is the ID of the bot.
	BotID *string `json:"bot_id"`
	// WorkflowID is the ID of the workflow.
	WorkflowID *string `json:"workflow_id"`
	// DeviceID is the ID of the device.
	DeviceID *int64 `json:"device_id"`
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
