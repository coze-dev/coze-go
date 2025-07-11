package coze

type websocketAudio struct {
	core           *core
	Speech         *websocketAudioSpeechBuild
	Transcriptions *websocketAudioTranscriptionBuild
}

func newWebsocketAudio(core *core) *websocketAudio {
	return &websocketAudio{
		core:           core,
		Speech:         newWebsocketAudioSpeechBuild(core),
		Transcriptions: newWebsocketAudioTranscriptionBuild(core),
	}
}
