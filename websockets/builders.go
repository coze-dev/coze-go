package websockets

// Speech returns a new speech client
func (a *AudioClient) Speech(opts ...SpeechClientOption) *SpeechClient {
	return NewSpeechClient(a.baseURL, a.auth, opts...)
}

// Transcriptions returns a new transcriptions client
func (a *AudioClient) Transcriptions(opts ...TranscriptionsClientOption) *TranscriptionsClient {
	return NewTranscriptionsClient(a.baseURL, a.auth, opts...)
}

// Create returns a new chat client
func (c *ChatClientBuilder) Create(opts ...ChatClientOption) *ChatClient {
	return NewChatClient(c.baseURL, c.auth, opts...)
}