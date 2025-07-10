package coze

// Speech returns a new speech client
func (a *websocketAudio) Speech(opts ...SpeechClientOption) *SpeechClient {
	return NewSpeechClient(a.baseURL, a.auth, opts...)
}

// Transcriptions returns a new transcriptions client
func (a *websocketAudio) Transcriptions() *TranscriptionsClient {
	return NewTranscriptionsClient(a.baseURL, a.auth, opts...)
}
