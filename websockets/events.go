package websockets

import (
	"encoding/json"
)

// WebSocketEvent represents a base WebSocket event
type WebSocketEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      json.RawMessage    `json:"data,omitempty"`
}

// EventDetail contains additional event details
type EventDetail struct {
	LogID         string `json:"logid,omitempty"`
	OriginMessage string `json:"origin_message,omitempty"`
}

// WebSocketErrorEvent represents an error event
type WebSocketErrorEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      *ErrorData         `json:"data,omitempty"`
}

// ErrorData contains error information
type ErrorData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// LimitConfig configures audio limits
type LimitConfig struct {
	Period       int `json:"period,omitempty"`
	MaxFrameNum  int `json:"max_frame_num,omitempty"`
}

// InputAudio configuration for audio input
type InputAudio struct {
	Format     string `json:"format,omitempty"`
	Codec      string `json:"codec,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
	Channel    int    `json:"channel,omitempty"`
	BitDepth   int    `json:"bit_depth,omitempty"`
}

// OpusConfig configures Opus audio output
type OpusConfig struct {
	Bitrate      int          `json:"bitrate,omitempty"`
	UseCBR       bool         `json:"use_cbr,omitempty"`
	FrameSizeMs  float64      `json:"frame_size_ms,omitempty"`
	LimitConfig  *LimitConfig `json:"limit_config,omitempty"`
}

// PCMConfig configures PCM audio output
type PCMConfig struct {
	SampleRate  int          `json:"sample_rate,omitempty"`
	FrameSizeMs float64      `json:"frame_size_ms,omitempty"`
	LimitConfig *LimitConfig `json:"limit_config,omitempty"`
}

// OutputAudio configuration for audio output
type OutputAudio struct {
	Codec      string      `json:"codec,omitempty"`
	PCMConfig  *PCMConfig  `json:"pcm_config,omitempty"`
	OpusConfig *OpusConfig `json:"opus_config,omitempty"`
	SpeechRate int         `json:"speech_rate,omitempty"`
	VoiceID    string      `json:"voice_id,omitempty"`
}

// Audio Speech Events

// SpeechUpdateEvent represents a speech update event
type SpeechUpdateEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      *SpeechUpdateData  `json:"data,omitempty"`
}

// SpeechUpdateData contains speech update configuration
type SpeechUpdateData struct {
	OutputAudio *OutputAudio `json:"output_audio,omitempty"`
}

// InputTextBufferAppendEvent represents text buffer append event
type InputTextBufferAppendEvent struct {
	EventType WebSocketEventType         `json:"event_type"`
	ID        string                     `json:"id,omitempty"`
	Detail    *EventDetail               `json:"detail,omitempty"`
	Data      *InputTextBufferAppendData `json:"data,omitempty"`
}

// InputTextBufferAppendData contains the text delta
type InputTextBufferAppendData struct {
	Delta string `json:"delta"`
}

// InputTextBufferCompleteEvent represents text buffer complete event
type InputTextBufferCompleteEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
}

// SpeechCreatedEvent represents speech created event
type SpeechCreatedEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      *SpeechCreatedData `json:"data,omitempty"`
}

// SpeechCreatedData contains speech session information
type SpeechCreatedData struct {
	SessionID string `json:"session_id"`
}

// SpeechAudioUpdateEvent represents speech audio update event
type SpeechAudioUpdateEvent struct {
	EventType WebSocketEventType     `json:"event_type"`
	ID        string                 `json:"id,omitempty"`
	Detail    *EventDetail           `json:"detail,omitempty"`
	Data      *SpeechAudioUpdateData `json:"data,omitempty"`
}

// SpeechAudioUpdateData contains audio delta
type SpeechAudioUpdateData struct {
	Delta string `json:"delta"` // Base64 encoded audio
}

// SpeechAudioCompletedEvent represents speech audio completed event
type SpeechAudioCompletedEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      *SpeechAudioCompletedData `json:"data,omitempty"`
}

// SpeechAudioCompletedData contains completion information
type SpeechAudioCompletedData struct {
	SessionID string `json:"session_id"`
}

// Audio Transcription Events

// TranscriptionsUpdateEvent represents transcriptions update event
type TranscriptionsUpdateEvent struct {
	EventType WebSocketEventType        `json:"event_type"`
	ID        string                    `json:"id,omitempty"`
	Detail    *EventDetail              `json:"detail,omitempty"`
	Data      *TranscriptionsUpdateData `json:"data,omitempty"`
}

// TranscriptionsUpdateData contains transcription configuration
type TranscriptionsUpdateData struct {
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

// InputAudioBufferAppendEvent represents audio buffer append event
type InputAudioBufferAppendEvent struct {
	EventType WebSocketEventType          `json:"event_type"`
	ID        string                      `json:"id,omitempty"`
	Detail    *EventDetail                `json:"detail,omitempty"`
	Data      *InputAudioBufferAppendData `json:"data,omitempty"`
}

// InputAudioBufferAppendData contains audio delta
type InputAudioBufferAppendData struct {
	Delta string `json:"delta"` // Base64 encoded audio
}

// InputAudioBufferCompleteEvent represents audio buffer complete event
type InputAudioBufferCompleteEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
}

// InputAudioBufferClearEvent represents audio buffer clear event
type InputAudioBufferClearEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
}

// TranscriptionsMessageUpdateEvent represents transcription message update event
type TranscriptionsMessageUpdateEvent struct {
	EventType WebSocketEventType               `json:"event_type"`
	ID        string                           `json:"id,omitempty"`
	Detail    *EventDetail                     `json:"detail,omitempty"`
	Data      *TranscriptionsMessageUpdateData `json:"data,omitempty"`
}

// TranscriptionsMessageUpdateData contains transcription content
type TranscriptionsMessageUpdateData struct {
	Content string `json:"content"`
}

// Chat Events

// ChatUpdateEvent represents chat update event
type ChatUpdateEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
	Data      *ChatUpdateData    `json:"data,omitempty"`
}

// ChatUpdateData contains chat configuration
type ChatUpdateData struct {
	BotID       string       `json:"bot_id,omitempty"`
	InputAudio  *InputAudio  `json:"input_audio,omitempty"`
	OutputAudio *OutputAudio `json:"output_audio,omitempty"`
}

// ConversationMessageCreateEvent represents conversation message create event
type ConversationMessageCreateEvent struct {
	EventType WebSocketEventType             `json:"event_type"`
	ID        string                         `json:"id,omitempty"`
	Detail    *EventDetail                   `json:"detail,omitempty"`
	Data      *ConversationMessageCreateData `json:"data,omitempty"`
}

// ConversationMessageCreateData contains message content
type ConversationMessageCreateData struct {
	Content string `json:"content"`
}

// ConversationMessageDeltaEvent represents conversation message delta event
type ConversationMessageDeltaEvent struct {
	EventType WebSocketEventType            `json:"event_type"`
	ID        string                        `json:"id,omitempty"`
	Detail    *EventDetail                  `json:"detail,omitempty"`
	Data      *ConversationMessageDeltaData `json:"data,omitempty"`
}

// ConversationMessageDeltaData contains message delta
type ConversationMessageDeltaData struct {
	Content string `json:"content"`
}

// ConversationAudioDeltaEvent represents conversation audio delta event
type ConversationAudioDeltaEvent struct {
	EventType WebSocketEventType          `json:"event_type"`
	ID        string                      `json:"id,omitempty"`
	Detail    *EventDetail                `json:"detail,omitempty"`
	Data      *ConversationAudioDeltaData `json:"data,omitempty"`
}

// ConversationAudioDeltaData contains audio delta
type ConversationAudioDeltaData struct {
	Content string `json:"content"` // Base64 encoded audio
}

// GetAudio returns the decoded audio bytes
func (c *ConversationAudioDeltaData) GetAudio() []byte {
	// In a real implementation, this would decode the base64 content
	// For now, we'll return the raw bytes
	return []byte(c.Content)
}

// ConversationChatCreatedEvent represents conversation chat created event
type ConversationChatCreatedEvent struct {
	EventType WebSocketEventType           `json:"event_type"`
	ID        string                       `json:"id,omitempty"`
	Detail    *EventDetail                 `json:"detail,omitempty"`
	Data      *ConversationChatCreatedData `json:"data,omitempty"`
}

// ConversationChatCreatedData contains chat information
type ConversationChatCreatedData struct {
	ChatID string `json:"chat_id"`
}

// ConversationChatCompletedEvent represents conversation chat completed event
type ConversationChatCompletedEvent struct {
	EventType WebSocketEventType             `json:"event_type"`
	ID        string                         `json:"id,omitempty"`
	Detail    *EventDetail                   `json:"detail,omitempty"`
	Data      *ConversationChatCompletedData `json:"data,omitempty"`
}

// ConversationChatCompletedData contains completion information
type ConversationChatCompletedData struct {
	ChatID string `json:"chat_id"`
}

// ConversationChatRequiresActionEvent represents conversation chat requires action event
type ConversationChatRequiresActionEvent struct {
	EventType WebSocketEventType                   `json:"event_type"`
	ID        string                               `json:"id,omitempty"`
	Detail    *EventDetail                         `json:"detail,omitempty"`
	Data      *ConversationChatRequiresActionData  `json:"data,omitempty"`
}

// ConversationChatRequiresActionData contains tool call requirements
type ConversationChatRequiresActionData struct {
	ChatID         string          `json:"chat_id"`
	RequiredAction *RequiredAction `json:"required_action,omitempty"`
}

// RequiredAction represents a required action
type RequiredAction struct {
	Type              string                     `json:"type"`
	SubmitToolOutputs *SubmitToolOutputsAction   `json:"submit_tool_outputs,omitempty"`
}

// SubmitToolOutputsAction represents tool outputs action
type SubmitToolOutputsAction struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

// ToolCall represents a tool call
type ToolCall struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}

// Function represents a function call
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ConversationChatSubmitToolOutputsEvent represents tool outputs submission event
type ConversationChatSubmitToolOutputsEvent struct {
	EventType WebSocketEventType                      `json:"event_type"`
	ID        string                                  `json:"id,omitempty"`
	Detail    *EventDetail                            `json:"detail,omitempty"`
	Data      *ConversationChatSubmitToolOutputsData  `json:"data,omitempty"`
}

// ConversationChatSubmitToolOutputsData contains tool outputs
type ConversationChatSubmitToolOutputsData struct {
	ChatID      string       `json:"chat_id"`
	ToolOutputs []ToolOutput `json:"tool_outputs"`
}

// ToolOutput represents a tool output
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// ConversationChatCancelEvent represents conversation chat cancel event
type ConversationChatCancelEvent struct {
	EventType WebSocketEventType          `json:"event_type"`
	ID        string                      `json:"id,omitempty"`
	Detail    *EventDetail                `json:"detail,omitempty"`
	Data      *ConversationChatCancelData `json:"data,omitempty"`
}

// ConversationChatCancelData contains cancel information
type ConversationChatCancelData struct {
	ChatID string `json:"chat_id"`
}

// ConversationChatCanceledEvent represents conversation chat canceled event
type ConversationChatCanceledEvent struct {
	EventType WebSocketEventType           `json:"event_type"`
	ID        string                       `json:"id,omitempty"`
	Detail    *EventDetail                 `json:"detail,omitempty"`
	Data      *ConversationChatCanceledData `json:"data,omitempty"`
}

// ConversationChatCanceledData contains cancellation information
type ConversationChatCanceledData struct {
	ChatID string `json:"chat_id"`
}