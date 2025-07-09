package coze

// Speech returns a new speech client
func (a *websocketAudio) Speech(opts ...SpeechClientOption) *SpeechClient {
	return NewSpeechClient(a.baseURL, a.auth, opts...)
}

// Transcriptions returns a new transcriptions client
func (a *websocketAudio) Transcriptions(opts ...TranscriptionsClientOption) *TranscriptionsClient {
	return NewTranscriptionsClient(a.baseURL, a.auth, opts...)
}

// Create returns a new chat client
func (c *websocketChatBuilder) Create(opts ...ChatClientOption) *ChatClient {
	return NewChatClient(c.baseURL, c.auth, opts...)
}
