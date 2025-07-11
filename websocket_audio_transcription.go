package coze

import "context"

func (r *websocketAudioTranscriptionBuild) Create(ctx context.Context, req *CreateWebsocketAudioTranscriptionReq) *WebSocketAudioTranscription {
	return newWebSocketAudioTranscriptionClient(ctx, r.core, req)
}

type CreateWebsocketAudioTranscriptionReq struct {
	// BotID is the ID of the bot.
	BotID *string `json:"bot_id"`
	// WorkflowID is the ID of the workflow.
	WorkflowID *string `json:"workflow_id"`
	// DeviceID is the ID of the device.
	DeviceID *int64 `json:"device_id"`
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
