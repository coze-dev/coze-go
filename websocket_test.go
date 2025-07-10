package coze

//
// import (
// 	"context"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )
//
// // MockAuth is a mock implementation of the Auth interface
// type MockAuth struct {
// 	mock.Mock
// }
//
// func (m *MockAuth) Token(ctx context.Context) string {
// 	args := m.Called()
// 	return args.String(0)
// }
//
// func TestWebSocketEventTypes(t *testing.T) {
// 	// Test that event types are correctly defined
// 	assert.Equal(t, "speech.created", string(WebSocketEventTypeSpeechCreated))
// 	assert.Equal(t, "speech.audio.update", string(WebSocketEventTypeSpeechAudioUpdate))
// 	assert.Equal(t, "conversation.message.delta", string(WebSocketEventTypeConversationMessageDelta))
// 	assert.Equal(t, "error", string(WebSocketEventTypeError))
// 	assert.Equal(t, "closed", string(WebSocketEventTypeClosed))
// }
//
// func TestSpeechClientCreation(t *testing.T) {
// 	mockAuth := &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	// Test speech client creation
// 	client := NewSpeechClient("wss://api.coze.com", mockAuth)
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.ws)
// }
//
// func TestTranscriptionsClientCreation(t *testing.T) {
// 	mockAuth := &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	// Test transcriptions client creation
// 	client := NewTranscriptionsClient("wss://api.coze.com", mockAuth)
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.ws)
// }
//
// func TestChatClientCreation(t *testing.T) {
// 	var mockAuth Auth = &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	// Test chat client creation
// 	client := NewChatClient("wss://api.coze.com", Auth(mockAuth), WithBotID("test-bot-id"))
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.ws)
// 	assert.Equal(t, "test-bot-id", client.botID)
// }
//
// func TestWebSocketsClientCreation(t *testing.T) {
// 	mockAuth := &MockAuth{}
//
// 	// Test websockets creation
// 	client := NewWebSockets("https://api.coze.com", mockAuth)
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.Audio)
// 	assert.NotNil(t, client.Chat)
// 	assert.Equal(t, "wss://api.coze.com", client.baseURL)
// }
//
// func TestEventHandlerRegistration(t *testing.T) {
// 	mockAuth := &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	speechClient := NewSpeechClient("wss://api.coze.com", mockAuth)
//
// 	called := false
// 	handler := &SpeechEventHandler{
// 		OnError: func(err error) error {
// 			called = true
// 			return nil
// 		},
// 	}
//
// 	// Register handlers
// 	handler.RegisterHandlers(speechClient)
//
// 	// Simulate an error event
// 	errorEvent := &WebSocketEvent{
// 		EventType: WebSocketEventTypeError,
// 		Data:      []byte(`{"error": "test error"}`),
// 	}
//
// 	// Get the registered handler and call it
// 	speechClient.ws.handlers[WebSocketEventTypeError](errorEvent)
//
// 	assert.True(t, called)
// }
//
// func TestAudioHelperFunctions(t *testing.T) {
// 	// Test base64 audio decoding
// 	testData := "SGVsbG8gV29ybGQ=" // "Hello World" in base64
// 	decoded, err := GetAudioFromDelta(testData)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Hello World", string(decoded))
// }
//
// func TestEventStructures(t *testing.T) {
// 	// Test WebSocketSpeechUpdateEventData structure
// 	outputAudio := &OutputAudio{
// 		Codec:      "pcm",
// 		SpeechRate: 0,
// 		VoiceID:    "test-voice",
// 	}
//
// 	speechUpdateData := &WebSocketSpeechUpdateEventData{
// 		OutputAudio: outputAudio,
// 	}
//
// 	assert.Equal(t, "pcm", speechUpdateData.OutputAudio.Codec)
// 	assert.Equal(t, "test-voice", speechUpdateData.OutputAudio.VoiceID)
//
// 	// Test WebSocketInputTextBufferAppendEventData structure
// 	textData := &WebSocketInputTextBufferAppendEventData{
// 		Delta: "Hello World",
// 	}
//
// 	assert.Equal(t, "Hello World", textData.Delta)
//
// 	// Test ConversationAudioDeltaData structure
// 	audioData := &ConversationAudioDeltaData{
// 		Content: "test-audio-content",
// 	}
//
// 	audio := audioData.GetAudio()
// 	assert.Equal(t, []byte("test-audio-content"), audio)
// }
//
// func TestChatClientMethods(t *testing.T) {
// 	mockAuth := &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	chatClient := NewChatClient("wss://api.coze.com", mockAuth, WithBotID("test-bot"))
//
// 	// Test configuration methods exist and don't panic
// 	assert.NotPanics(t, func() {
// 		// These would normally send events, but we're just testing the method exists
// 		chatClient.UpdateChat("test-bot", nil, nil)
// 		chatClient.CreateMessage("test message")
// 		chatClient.ClearConversation()
// 		chatClient.CancelChat("test-chat-id")
// 	})
// }
//
// func TestWebSocketClientOptions(t *testing.T) {
// 	mockAuth := &MockAuth{}
// 	mockAuth.On("Token").Return("Bearer test-token", nil)
//
// 	// Test speech client with options
// 	speechClient := NewSpeechClient("wss://api.coze.com", mockAuth,
// 		WithOutputAudio(&OutputAudio{Codec: "pcm"}))
// 	assert.NotNil(t, speechClient)
//
// 	// Test transcriptions client with options
// 	transcriptionsClient := NewTranscriptionsClient("wss://api.coze.com", mockAuth,
// 		WithInputAudio(&InputAudio{Format: "wav"}))
// 	assert.NotNil(t, transcriptionsClient)
//
// 	// Test chat client with options
// 	chatClient := NewChatClient("wss://api.coze.com", mockAuth,
// 		WithBotID("test-bot"),
// 		WithChatInputAudio(&InputAudio{Format: "wav"}),
// 		WithChatOutputAudio(&OutputAudio{Codec: "pcm"}))
// 	assert.NotNil(t, chatClient)
// 	assert.Equal(t, "test-bot", chatClient.botID)
// }
//
// func TestURLConversion(t *testing.T) {
// 	mockAuth := &MockAuth{}
//
// 	// Test HTTP to WebSocket URL conversion
// 	client := NewWebSockets("http://api.coze.com", mockAuth)
// 	assert.Equal(t, "ws://api.coze.com", client.baseURL)
//
// 	// Test HTTPS to WebSocket URL conversion
// 	client = NewWebSockets("https://api.coze.com", mockAuth)
// 	assert.Equal(t, "wss://api.coze.com", client.baseURL)
//
// 	// Test WebSocket URL remains unchanged
// 	client = NewWebSockets("wss://api.coze.com", mockAuth)
// 	assert.Equal(t, "wss://api.coze.com", client.baseURL)
// }
