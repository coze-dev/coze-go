package coze

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// WebSocketEventType websocket 事件类型
type WebSocketEventType string

const (
	// common

	WebSocketEventTypeClientError WebSocketEventType = "client_error" // sdk error
	WebSocketEventTypeClosed      WebSocketEventType = "closed"       // connection closed
	WebSocketEventTypeError       WebSocketEventType = "error"        // 发生错误

	// v1/audio/speech

	// req

	WebSocketEventTypeSpeechUpdate            WebSocketEventType = "speech.update"              // 更新语音合成配置
	WebSocketEventTypeInputTextBufferAppend   WebSocketEventType = "input_text_buffer.append"   // 流式输入文字
	WebSocketEventTypeInputTextBufferComplete WebSocketEventType = "input_text_buffer.complete" // 提交文字

	// resp

	WebSocketEventTypeSpeechCreated            WebSocketEventType = "speech.created"              // 语音合成连接成功
	WebSocketEventTypeSpeechUpdated            WebSocketEventType = "speech.updated"              // 配置更新完成
	WebSocketEventTypeInputTextBufferCompleted WebSocketEventType = "input_text_buffer.completed" // input_text_buffer 提交完成
	WebSocketEventTypeSpeechAudioUpdate        WebSocketEventType = "speech.audio.update"         // 合成增量语音
	WebSocketEventTypeSpeechAudioCompleted     WebSocketEventType = "speech.audio.completed"      // 合成完成

	// v1/audio/transcriptions

	// req

	WebSocketEventTypeTranscriptionsUpdate     WebSocketEventType = "transcriptions.update"       // 更新语音识别配置
	WebSocketEventTypeInputAudioBufferAppend   WebSocketEventType = "input_audio_buffer.append"   // 流式上传音频片段
	WebSocketEventTypeInputAudioBufferComplete WebSocketEventType = "input_audio_buffer.complete" // 提交音频
	WebSocketEventTypeInputAudioBufferClear    WebSocketEventType = "input_audio_buffer.clear"    // 清除缓冲区音频

	// resp

	WebSocketEventTypeTranscriptionsCreated          WebSocketEventType = "transcriptions.created"           // 连接成功
	WebSocketEventTypeTranscriptionsUpdated          WebSocketEventType = "transcriptions.updated"           // 配置更新成功
	WebSocketEventTypeInputAudioBufferCompleted      WebSocketEventType = "input_audio_buffer.completed"     // 音频提交完成
	WebSocketEventTypeInputAudioBufferCleared        WebSocketEventType = "input_audio_buffer.cleared"       // 音频清除成功
	WebSocketEventTypeTranscriptionsMessageUpdate    WebSocketEventType = "transcriptions.message.update"    // 识别出文字
	WebSocketEventTypeTranscriptionsMessageCompleted WebSocketEventType = "transcriptions.message.completed" // 识别完成

	// v1/chat

	// req

	WebSocketEventTypeChatUpdate WebSocketEventType = "chat.update" // 更新对话配置
	// WebSocketEventTypeInputAudioBufferAppend            WebSocketEventType = "input_audio_buffer.append"             // 流式上传音频片段
	// WebSocketEventTypeInputAudioBufferComplete          WebSocketEventType = "input_audio_buffer.complete"           // 提交音频
	// WebSocketEventTypeInputAudioBufferClear             WebSocketEventType = "input_audio_buffer.clear"              // 清除缓冲区音频
	WebSocketEventTypeConversationMessageCreate         WebSocketEventType = "conversation.message.create"           // 手动提交对话内容
	WebSocketEventTypeConversationClear                 WebSocketEventType = "conversation.clear"                    // 清除上下文
	WebSocketEventTypeConversationChatSubmitToolOutputs WebSocketEventType = "conversation.chat.submit_tool_outputs" // 提交端插件执行结果
	WebSocketEventTypeConversationChatCancel            WebSocketEventType = "conversation.chat.cancel"              // 打断智能体输出

	// resp

	WebSocketEventTypeChatCreated                    WebSocketEventType = "chat.created"                      // 对话连接成功
	WebSocketEventTypeChatUpdated                    WebSocketEventType = "chat.updated"                      // 对话配置成功
	WebSocketEventTypeConversationChatCreated        WebSocketEventType = "conversation.chat.created"         // 对话开始
	WebSocketEventTypeConversationChatInProgress     WebSocketEventType = "conversation.chat.in_progress"     // 对话正在处理
	WebSocketEventTypeConversationMessageDelta       WebSocketEventType = "conversation.message.delta"        // 增量消息
	WebSocketEventTypeConversationAudioSentenceStart WebSocketEventType = "conversation.audio.sentence_start" // 增量语音字幕
	WebSocketEventTypeConversationAudioDelta         WebSocketEventType = "conversation.audio.delta"          // 增量语音
	WebSocketEventTypeConversationMessageCompleted   WebSocketEventType = "conversation.message.completed"    // 消息完成
	WebSocketEventTypeConversationAudioCompleted     WebSocketEventType = "conversation.audio.completed"      // 语音回复完成
	WebSocketEventTypeConversationChatCompleted      WebSocketEventType = "conversation.chat.completed"       // 对话完成
	WebSocketEventTypeConversationChatFailed         WebSocketEventType = "conversation.chat.failed"          // 对话失败
	// WebSocketEventTypeInputAudioBufferCompleted            WebSocketEventType = "input_audio_buffer.completed"            // 音频提交完成
	// WebSocketEventTypeInputAudioBufferCleared              WebSocketEventType = "input_audio_buffer.cleared"              // 音频清除成功
	WebSocketEventTypeConversationCleared                  WebSocketEventType = "conversation.cleared"                    // 上下文清除完成
	WebSocketEventTypeConversationChatCanceled             WebSocketEventType = "conversation.chat.canceled"              // 智能体输出中断
	WebSocketEventTypeConversationAudioTranscriptUpdate    WebSocketEventType = "conversation.audio_transcript.update"    // 用户语音识别字幕
	WebSocketEventTypeConversationAudioTranscriptCompleted WebSocketEventType = "conversation.audio_transcript.completed" // 用户语音识别完成
	WebSocketEventTypeConversationChatRequiresAction       WebSocketEventType = "conversation.chat.requires_action"       // 端插件请求
	WebSocketEventTypeInputAudioBufferSpeechStarted        WebSocketEventType = "input_audio_buffer.speech_started"       // 用户开始说话
	WebSocketEventTypeInputAudioBufferSpeechStopped        WebSocketEventType = "input_audio_buffer.speech_stopped"       // 用户结束说话
)

var websocketEvents = map[string]reflect.Type{
	// common
	string(WebSocketEventTypeClientError): reflect.TypeOf(WebSocketClientErrorEvent{}),
	string(WebSocketEventTypeClosed):      reflect.TypeOf(WebSocketClosedEvent{}),
	string(WebSocketEventTypeError):       reflect.TypeOf(WebSocketErrorEvent{}),

	// v1/audio/speech req
	string(WebSocketEventTypeSpeechUpdate):            reflect.TypeOf(WebSocketSpeechUpdateEvent{}),
	string(WebSocketEventTypeInputTextBufferAppend):   reflect.TypeOf(WebSocketInputTextBufferAppendEvent{}),
	string(WebSocketEventTypeInputTextBufferComplete): reflect.TypeOf(WebSocketInputTextBufferCompleteEvent{}),
	// v1/audio/speech resp
	string(WebSocketEventTypeSpeechCreated):            reflect.TypeOf(WebSocketSpeechCreatedEvent{}),
	string(WebSocketEventTypeSpeechUpdated):            reflect.TypeOf(WebSocketSpeechUpdatedEvent{}),
	string(WebSocketEventTypeInputTextBufferCompleted): reflect.TypeOf(WebSocketInputTextBufferCompletedEvent{}),
	string(WebSocketEventTypeSpeechAudioUpdate):        reflect.TypeOf(WebSocketSpeechAudioUpdateEvent{}),
	string(WebSocketEventTypeSpeechAudioCompleted):     reflect.TypeOf(WebSocketSpeechAudioCompletedEvent{}),
	// v1/audio/transcriptions req
	string(WebSocketEventTypeTranscriptionsUpdate):     reflect.TypeOf(WebSocketTranscriptionsUpdateEvent{}),
	string(WebSocketEventTypeInputAudioBufferAppend):   reflect.TypeOf(WebSocketInputAudioBufferAppendEvent{}),
	string(WebSocketEventTypeInputAudioBufferComplete): reflect.TypeOf(WebSocketInputAudioBufferCompleteEvent{}),
	string(WebSocketEventTypeInputAudioBufferClear):    reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
	// v1/audio/transcriptions resp
	string(WebSocketEventTypeTranscriptionsCreated):          reflect.TypeOf(WebSocketTranscriptionsCreatedEvent{}),
	string(WebSocketEventTypeTranscriptionsUpdated):          reflect.TypeOf(WebSocketTranscriptionsUpdatedEvent{}),
	string(WebSocketEventTypeInputAudioBufferCompleted):      reflect.TypeOf(WebSocketInputAudioBufferCompletedEvent{}),
	string(WebSocketEventTypeInputAudioBufferCleared):        reflect.TypeOf(WebSocketInputAudioBufferClearedEvent{}),
	string(WebSocketEventTypeTranscriptionsMessageUpdate):    reflect.TypeOf(WebSocketTranscriptionsMessageUpdateEvent{}),
	string(WebSocketEventTypeTranscriptionsMessageCompleted): reflect.TypeOf(WebSocketTranscriptionsMessageCompletedEvent{}),
	// v1/chat req
	string(WebSocketEventTypeChatUpdate): reflect.TypeOf(WebSocketChatUpdateEvent{}),
	// string(WebSocketEventTypeInputAudioBufferAppend):   reflect.TypeOf(WebSocketInputAudioBufferAppendEvent{}),
	// string(WebSocketEventTypeInputAudioBufferComplete): reflect.TypeOf(WebSocketInputAudioBufferCompleteEvent{}),
	// string(WebSocketEventTypeInputAudioBufferClear):    reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
	string(WebSocketEventTypeConversationMessageCreate):         reflect.TypeOf(WebSocketConversationMessageCreateEvent{}),
	string(WebSocketEventTypeConversationClear):                 reflect.TypeOf(WebSocketConversationClearEvent{}),
	string(WebSocketEventTypeConversationChatSubmitToolOutputs): reflect.TypeOf(WebSocketConversationChatSubmitToolOutputsEvent{}),
	string(WebSocketEventTypeConversationChatCancel):            reflect.TypeOf(WebSocketConversationChatCancelEvent{}),
	// v1/chat resp
	string(WebSocketEventTypeChatCreated):                    reflect.TypeOf(WebSocketChatCreatedEvent{}),
	string(WebSocketEventTypeChatUpdated):                    reflect.TypeOf(WebSocketChatUpdatedEvent{}),
	string(WebSocketEventTypeConversationChatCreated):        reflect.TypeOf(WebSocketConversationChatCreatedEvent{}),
	string(WebSocketEventTypeConversationChatInProgress):     reflect.TypeOf(WebSocketConversationChatInProgressEvent{}),
	string(WebSocketEventTypeConversationMessageDelta):       reflect.TypeOf(WebSocketConversationMessageDeltaEvent{}),
	string(WebSocketEventTypeConversationAudioSentenceStart): reflect.TypeOf(WebSocketConversationAudioSentenceStartEvent{}),
	string(WebSocketEventTypeConversationAudioDelta):         reflect.TypeOf(WebSocketConversationAudioDeltaEvent{}),
	string(WebSocketEventTypeConversationMessageCompleted):   reflect.TypeOf(WebSocketConversationMessageCompletedEvent{}),
	string(WebSocketEventTypeConversationAudioCompleted):     reflect.TypeOf(WebSocketConversationAudioCompletedEvent{}),
	string(WebSocketEventTypeConversationChatCompleted):      reflect.TypeOf(WebSocketConversationChatCompletedEvent{}),
	string(WebSocketEventTypeConversationChatFailed):         reflect.TypeOf(WebSocketConversationChatFailedEvent{}),
	// string(WebSocketEventTypeInputAudioBufferCompleted):      reflect.TypeOf(WebSocketInputAudioBufferCompletedEvent{}),
	// string(WebSocketEventTypeInputAudioBufferCleared):        reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
	string(WebSocketEventTypeConversationCleared):                  reflect.TypeOf(WebSocketConversationClearedEvent{}),
	string(WebSocketEventTypeConversationChatCanceled):             reflect.TypeOf(WebSocketConversationChatCanceledEvent{}),
	string(WebSocketEventTypeConversationAudioTranscriptUpdate):    reflect.TypeOf(WebSocketConversationAudioTranscriptUpdateEvent{}),
	string(WebSocketEventTypeConversationAudioTranscriptCompleted): reflect.TypeOf(WebSocketConversationAudioTranscriptCompletedEvent{}),
	string(WebSocketEventTypeConversationChatRequiresAction):       reflect.TypeOf(WebSocketConversationChatRequiresActionEvent{}),
	string(WebSocketEventTypeInputAudioBufferSpeechStarted):        reflect.TypeOf(WebSocketInputAudioBufferSpeechStartedEvent{}),
	string(WebSocketEventTypeInputAudioBufferSpeechStopped):        reflect.TypeOf(WebSocketInputAudioBufferSpeechStoppedEvent{}),
}

var audioSpeechResponseEventTypes = []WebSocketEventType{
	WebSocketEventTypeClientError,
	WebSocketEventTypeClosed,
	WebSocketEventTypeError,

	WebSocketEventTypeSpeechCreated,
	WebSocketEventTypeSpeechUpdated,
	WebSocketEventTypeInputTextBufferCompleted,
	WebSocketEventTypeSpeechAudioUpdate,
	WebSocketEventTypeSpeechAudioCompleted,
}

var audioTranscriptionResponseEventTypes = []WebSocketEventType{
	WebSocketEventTypeClientError,
	WebSocketEventTypeClosed,
	WebSocketEventTypeError,

	WebSocketEventTypeTranscriptionsCreated,
	WebSocketEventTypeTranscriptionsUpdated,
	WebSocketEventTypeInputAudioBufferCompleted,
	WebSocketEventTypeInputAudioBufferCleared,
	WebSocketEventTypeTranscriptionsMessageUpdate,
	WebSocketEventTypeTranscriptionsMessageCompleted,
}

var chatResponseEventTypes = []WebSocketEventType{
	WebSocketEventTypeClientError,
	WebSocketEventTypeClosed,
	WebSocketEventTypeError,

	WebSocketEventTypeChatCreated,
	WebSocketEventTypeChatUpdated,
	WebSocketEventTypeConversationChatCreated,
	WebSocketEventTypeConversationChatInProgress,
	WebSocketEventTypeConversationMessageDelta,
	WebSocketEventTypeConversationAudioSentenceStart,
	WebSocketEventTypeConversationAudioDelta,
	WebSocketEventTypeConversationMessageCompleted,
	WebSocketEventTypeConversationAudioCompleted,
	WebSocketEventTypeConversationChatCompleted,
	WebSocketEventTypeConversationChatFailed,
	WebSocketEventTypeInputAudioBufferCompleted,
	WebSocketEventTypeInputAudioBufferCleared,
	WebSocketEventTypeConversationCleared,
	WebSocketEventTypeConversationChatCanceled,
	WebSocketEventTypeConversationAudioTranscriptUpdate,
	WebSocketEventTypeConversationAudioTranscriptCompleted,
	WebSocketEventTypeConversationChatRequiresAction,
	WebSocketEventTypeInputAudioBufferSpeechStarted,
	WebSocketEventTypeInputAudioBufferSpeechStopped,
}

const websocketEventTypeSize = 44

var websocketEventTypes = [websocketEventTypeSize]WebSocketEventType{
	WebSocketEventTypeClientError,
	WebSocketEventTypeClosed,
	WebSocketEventTypeError,
	WebSocketEventTypeSpeechUpdate,
	WebSocketEventTypeInputTextBufferAppend,
	WebSocketEventTypeInputTextBufferComplete,
	WebSocketEventTypeSpeechCreated,
	WebSocketEventTypeSpeechUpdated,
	WebSocketEventTypeInputTextBufferCompleted,
	WebSocketEventTypeSpeechAudioUpdate,
	WebSocketEventTypeSpeechAudioCompleted,
	WebSocketEventTypeTranscriptionsUpdate,
	WebSocketEventTypeInputAudioBufferAppend,
	WebSocketEventTypeInputAudioBufferComplete,
	WebSocketEventTypeInputAudioBufferClear,
	WebSocketEventTypeTranscriptionsCreated,
	WebSocketEventTypeTranscriptionsUpdated,
	WebSocketEventTypeInputAudioBufferCompleted,
	WebSocketEventTypeInputAudioBufferCleared,
	WebSocketEventTypeTranscriptionsMessageUpdate,
	WebSocketEventTypeTranscriptionsMessageCompleted,
	WebSocketEventTypeChatUpdate,
	WebSocketEventTypeConversationMessageCreate,
	WebSocketEventTypeConversationClear,
	WebSocketEventTypeConversationChatSubmitToolOutputs,
	WebSocketEventTypeConversationChatCancel,
	WebSocketEventTypeChatCreated,
	WebSocketEventTypeChatUpdated,
	WebSocketEventTypeConversationChatCreated,
	WebSocketEventTypeConversationChatInProgress,
	WebSocketEventTypeConversationMessageDelta,
	WebSocketEventTypeConversationAudioSentenceStart,
	WebSocketEventTypeConversationAudioDelta,
	WebSocketEventTypeConversationMessageCompleted,
	WebSocketEventTypeConversationAudioCompleted,
	WebSocketEventTypeConversationChatCompleted,
	WebSocketEventTypeConversationChatFailed,
	WebSocketEventTypeConversationCleared,
	WebSocketEventTypeConversationChatCanceled,
	WebSocketEventTypeConversationAudioTranscriptUpdate,
	WebSocketEventTypeConversationAudioTranscriptCompleted,
	WebSocketEventTypeConversationChatRequiresAction,
	WebSocketEventTypeInputAudioBufferSpeechStarted,
	WebSocketEventTypeInputAudioBufferSpeechStopped,
}

var websocketEventTypeIndex = map[WebSocketEventType]int{}

func init() {
	for k, v := range websocketEventTypes {
		websocketEventTypeIndex[v] = k
	}
}

type commonWebSocketEvent struct {
	baseWebSocketEvent
	Data json.RawMessage `json:"data,omitempty"`
}

func parseWebSocketEvent(message []byte) (IWebSocketEvent, error) {
	var common commonWebSocketEvent
	if err := json.Unmarshal(message, &common); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	eventTypeRef, ok := websocketEvents[string(common.GetEventType())]
	if !ok {
		return &common, nil
	}

	event := reflect.New(eventTypeRef).Interface()
	if err := json.Unmarshal(message, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return any(event).(IWebSocketEvent), nil
}

type baseWebSocketEvent struct {
	EventType WebSocketEventType `json:"event_type"`
	ID        string             `json:"id,omitempty"`
	Detail    *EventDetail       `json:"detail,omitempty"`
}

func (r baseWebSocketEvent) GetEventType() WebSocketEventType {
	return r.EventType
}

func (r baseWebSocketEvent) GetID() string {
	return r.ID
}

func (r baseWebSocketEvent) GetDetail() *EventDetail {
	return r.Detail
}

// EventDetail contains additional event details
type EventDetail struct {
	LogID         string `json:"logid,omitempty"`
	RespondAt     string `json:"respond_at,omitempty"`
	OriginMessage string `json:"origin_message,omitempty"`
}
