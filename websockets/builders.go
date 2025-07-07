package websockets

import (
	"github.com/coze-dev/coze-go/audio"
	"github.com/coze-dev/coze-go/chat"
)

// Speech returns a new speech client
func (a *AudioClient) Speech(opts ...audio.SpeechClientOption) *audio.SpeechClient {
	return audio.NewSpeechClient(a.baseURL, a.auth, opts...)
}

// Transcriptions returns a new transcriptions client
func (a *AudioClient) Transcriptions(opts ...audio.TranscriptionsClientOption) *audio.TranscriptionsClient {
	return audio.NewTranscriptionsClient(a.baseURL, a.auth, opts...)
}

// Create returns a new chat client
func (c *ChatClientBuilder) Create(opts ...chat.ChatClientOption) *chat.ChatClient {
	return chat.NewChatClient(c.baseURL, c.auth, opts...)
}