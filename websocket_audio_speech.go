package coze

// Speech returns a new speech client
func (a *websocketAudio) Speech(opts ...SpeechClientOption) *WebSocketAudioSpeech {
	return NewSpeechClient(a.baseURL, a.auth, opts...)
}
