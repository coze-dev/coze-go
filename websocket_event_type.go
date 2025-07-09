package coze

// WebSocketEventType represents different types of WebSocket events
type WebSocketEventType string

// Common WebSocket events
const (
	// common events

	EventTypeClientError WebSocketEventType = "client_error"
	EventTypeClosed      WebSocketEventType = "closed"
	EventTypeError       WebSocketEventType = "error"

	// audio speech events

	EventTypeSpeechUpdate             WebSocketEventType = "speech.update"
	EventTypeInputTextBufferAppend    WebSocketEventType = "input_text_buffer.append"
	EventTypeInputTextBufferComplete  WebSocketEventType = "input_text_buffer.complete"
	EventTypeSpeechCreated            WebSocketEventType = "speech.created"
	EventTypeSpeechUpdated            WebSocketEventType = "speech.updated"
	EventTypeInputTextBufferCompleted WebSocketEventType = "input_text_buffer.completed"
	EventTypeSpeechAudioUpdate        WebSocketEventType = "speech.audio.update"
	EventTypeSpeechAudioCompleted     WebSocketEventType = "speech.audio.completed"

	// audio transcription events

	EventTypeTranscriptionsUpdate           WebSocketEventType = "transcriptions.update"
	EventTypeInputAudioBufferAppend         WebSocketEventType = "input_audio_buffer.append"
	EventTypeInputAudioBufferComplete       WebSocketEventType = "input_audio_buffer.complete"
	EventTypeInputAudioBufferClear          WebSocketEventType = "input_audio_buffer.clear"
	EventTypeTranscriptionsCreated          WebSocketEventType = "transcriptions.created"
	EventTypeTranscriptionsUpdated          WebSocketEventType = "transcriptions.updated"
	EventTypeInputAudioBufferCompleted      WebSocketEventType = "input_audio_buffer.completed"
	EventTypeInputAudioBufferCleared        WebSocketEventType = "input_audio_buffer.cleared"
	EventTypeTranscriptionsMessageUpdate    WebSocketEventType = "transcriptions.message.update"
	EventTypeTranscriptionsMessageCompleted WebSocketEventType = "transcriptions.message.completed"

	// chat events

	EventTypeChatUpdate                           WebSocketEventType = "chat.update"
	EventTypeConversationMessageCreate            WebSocketEventType = "conversation.message.create"
	EventTypeConversationClear                    WebSocketEventType = "conversation.clear"
	EventTypeConversationChatSubmitToolOutputs    WebSocketEventType = "conversation.chat.submit_tool_outputs"
	EventTypeConversationChatCancel               WebSocketEventType = "conversation.chat.cancel"
	EventTypeChatCreated                          WebSocketEventType = "chat.created"
	EventTypeChatUpdated                          WebSocketEventType = "chat.updated"
	EventTypeConversationChatCreated              WebSocketEventType = "conversation.chat.created"
	EventTypeConversationChatInProgress           WebSocketEventType = "conversation.chat.in_progress"
	EventTypeConversationMessageDelta             WebSocketEventType = "conversation.message.delta"
	EventTypeConversationAudioSentenceStart       WebSocketEventType = "conversation.audio.sentence_start"
	EventTypeConversationAudioDelta               WebSocketEventType = "conversation.audio.delta"
	EventTypeConversationMessageCompleted         WebSocketEventType = "conversation.message.completed"
	EventTypeConversationAudioCompleted           WebSocketEventType = "conversation.audio.completed"
	EventTypeConversationChatCompleted            WebSocketEventType = "conversation.chat.completed"
	EventTypeConversationChatFailed               WebSocketEventType = "conversation.chat.failed"
	EventTypeConversationCleared                  WebSocketEventType = "conversation.cleared"
	EventTypeConversationChatCanceled             WebSocketEventType = "conversation.chat.canceled"
	EventTypeConversationAudioTranscriptUpdate    WebSocketEventType = "conversation.audio_transcript.update"
	EventTypeConversationAudioTranscriptCompleted WebSocketEventType = "conversation.audio_transcript.completed"
	EventTypeConversationChatRequiresAction       WebSocketEventType = "conversation.chat.requires_action"
	EventTypeInputAudioBufferSpeechStarted        WebSocketEventType = "input_audio_buffer.speech_started"
	EventTypeInputAudioBufferSpeechStopped        WebSocketEventType = "input_audio_buffer.speech_stopped"
)
