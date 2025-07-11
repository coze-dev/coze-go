package coze

// Transcriptions returns a new transcriptions client
func (a *websocketAudio) Transcriptions() *WebSocketAudioTranscription {
	return newWebSocketAudioTranscriptionClient()
}
