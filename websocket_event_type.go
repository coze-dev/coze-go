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

	EventTypeClientError WebSocketEventType = "client_error" // sdk error
	EventTypeClosed      WebSocketEventType = "closed"       // connection closed
	EventTypeError       WebSocketEventType = "error"        // 发生错误

	// v1/audio/speech

	// req

	EventTypeSpeechUpdate            WebSocketEventType = "speech.update"              // 更新语音合成配置
	EventTypeInputTextBufferAppend   WebSocketEventType = "input_text_buffer.append"   // 流式输入文字
	EventTypeInputTextBufferComplete WebSocketEventType = "input_text_buffer.complete" // 提交文字

	// resp

	EventTypeSpeechCreated            WebSocketEventType = "speech.created"              // 语音合成连接成功
	EventTypeSpeechUpdated            WebSocketEventType = "speech.updated"              // 配置更新完成
	EventTypeInputTextBufferCompleted WebSocketEventType = "input_text_buffer.completed" // input_text_buffer 提交完成
	EventTypeSpeechAudioUpdate        WebSocketEventType = "speech.audio.update"         // 合成增量语音
	EventTypeSpeechAudioCompleted     WebSocketEventType = "speech.audio.completed"      // 合成完成

	// v1/audio/transcriptions

	// req

	EventTypeTranscriptionsUpdate     WebSocketEventType = "transcriptions.update"       // 更新语音识别配置
	EventTypeInputAudioBufferAppend   WebSocketEventType = "input_audio_buffer.append"   // 流式上传音频片段
	EventTypeInputAudioBufferComplete WebSocketEventType = "input_audio_buffer.complete" // 提交音频
	EventTypeInputAudioBufferClear    WebSocketEventType = "input_audio_buffer.clear"    // 清除缓冲区音频

	// resp

	EventTypeTranscriptionsCreated          WebSocketEventType = "transcriptions.created"           // 连接成功
	EventTypeTranscriptionsUpdated          WebSocketEventType = "transcriptions.updated"           // 配置更新成功
	EventTypeInputAudioBufferCompleted      WebSocketEventType = "input_audio_buffer.completed"     // 音频提交完成
	EventTypeInputAudioBufferCleared        WebSocketEventType = "input_audio_buffer.cleared"       // 音频清除成功
	EventTypeTranscriptionsMessageUpdate    WebSocketEventType = "transcriptions.message.update"    // 识别出文字
	EventTypeTranscriptionsMessageCompleted WebSocketEventType = "transcriptions.message.completed" // 识别完成

	// v1/chat

	// req

	EventTypeChatUpdate WebSocketEventType = "chat.update" // 更新对话配置
	// EventTypeInputAudioBufferAppend   WebSocketEventType = "input_audio_buffer.append"   // 流式上传音频片段
	// EventTypeInputAudioBufferComplete WebSocketEventType = "input_audio_buffer.complete" // 提交音频
	// EventTypeInputAudioBufferClear    WebSocketEventType = "input_audio_buffer.clear"    // 清除缓冲区音频
	EventTypeConversationMessageCreate         WebSocketEventType = "conversation.message.create"           // 手动提交对话内容
	EventTypeConversationClear                 WebSocketEventType = "conversation.clear"                    // 清除上下文
	EventTypeConversationChatSubmitToolOutputs WebSocketEventType = "conversation.chat.submit_tool_outputs" // 提交端插件执行结果
	EventTypeConversationChatCancel            WebSocketEventType = "conversation.chat.cancel"              // 打断智能体输出

	// resp

	EventTypeChatCreated                    WebSocketEventType = "chat.created"                      // 对话连接成功
	EventTypeChatUpdated                    WebSocketEventType = "chat.updated"                      // 对话配置成功
	EventTypeConversationChatCreated        WebSocketEventType = "conversation.chat.created"         // 对话开始
	EventTypeConversationChatInProgress     WebSocketEventType = "conversation.chat.in_progress"     // 对话正在处理
	EventTypeConversationMessageDelta       WebSocketEventType = "conversation.message.delta"        // 增量消息
	EventTypeConversationAudioSentenceStart WebSocketEventType = "conversation.audio.sentence_start" // 增量语音字幕
	EventTypeConversationAudioDelta         WebSocketEventType = "conversation.audio.delta"          // 增量语音
	EventTypeConversationMessageCompleted   WebSocketEventType = "conversation.message.completed"    // 消息完成
	EventTypeConversationAudioCompleted     WebSocketEventType = "conversation.audio.completed"      // 语音回复完成
	EventTypeConversationChatCompleted      WebSocketEventType = "conversation.chat.completed"       // 对话完成
	EventTypeConversationChatFailed         WebSocketEventType = "conversation.chat.failed"          // 对话失败
	// EventTypeInputAudioBufferCompleted            WebSocketEventType = "input_audio_buffer.completed"            // 音频提交完成
	// EventTypeInputAudioBufferCleared              WebSocketEventType = "input_audio_buffer.cleared"              // 音频清除成功
	EventTypeConversationCleared                  WebSocketEventType = "conversation.cleared"                    // 上下文清除完成
	EventTypeConversationChatCanceled             WebSocketEventType = "conversation.chat.canceled"              // 智能体输出中断
	EventTypeConversationAudioTranscriptUpdate    WebSocketEventType = "conversation.audio_transcript.update"    // 用户语音识别字幕
	EventTypeConversationAudioTranscriptCompleted WebSocketEventType = "conversation.audio_transcript.completed" // 用户语音识别完成
	EventTypeConversationChatRequiresAction       WebSocketEventType = "conversation.chat.requires_action"       // 端插件请求
	EventTypeInputAudioBufferSpeechStarted        WebSocketEventType = "input_audio_buffer.speech_started"       // 用户开始说话
	EventTypeInputAudioBufferSpeechStopped        WebSocketEventType = "input_audio_buffer.speech_stopped"       // 用户结束说话
)

var websocketEvents = map[string]reflect.Type{}

func registerWebSocketEvent() {
	websocketEvents = map[string]reflect.Type{
		// common
		// string(EventTypeClientError): reflect.TypeOf(webclient{}),
		// string(EventTypeClosed): reflect.TypeOf(webclo{}),
		string(EventTypeError): reflect.TypeOf(WebSocketErrorEvent{}),

		// v1/audio/speech req
		string(EventTypeSpeechUpdate):            reflect.TypeOf(WebSocketSpeechUpdateEvent{}),
		string(EventTypeInputTextBufferAppend):   reflect.TypeOf(WebSocketInputTextBufferAppendEvent{}),
		string(EventTypeInputTextBufferComplete): reflect.TypeOf(WebSocketInputTextBufferCompleteEvent{}),
		// v1/audio/speech resp
		string(EventTypeSpeechCreated):            reflect.TypeOf(WebSocketSpeechCreatedEvent{}),
		string(EventTypeSpeechUpdated):            reflect.TypeOf(WebSocketSpeechUpdatedEvent{}),
		string(EventTypeInputTextBufferCompleted): reflect.TypeOf(WebSocketInputTextBufferCompleteEvent{}),
		string(EventTypeSpeechAudioUpdate):        reflect.TypeOf(WebSocketSpeechAudioUpdateEvent{}),
		string(EventTypeSpeechAudioCompleted):     reflect.TypeOf(WebSocketSpeechAudioCompletedEvent{}),
		// v1/audio/transcriptions req
		string(EventTypeTranscriptionsUpdate):     reflect.TypeOf(WebSocketTranscriptionsUpdateEvent{}),
		string(EventTypeInputAudioBufferAppend):   reflect.TypeOf(WebSocketInputAudioBufferAppendEvent{}),
		string(EventTypeInputAudioBufferComplete): reflect.TypeOf(WebSocketInputAudioBufferCompleteEvent{}),
		string(EventTypeInputAudioBufferClear):    reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
		// v1/audio/transcriptions resp
		string(EventTypeTranscriptionsCreated):          reflect.TypeOf(WebSocketTranscriptionsCreatedEvent{}),
		string(EventTypeTranscriptionsUpdated):          reflect.TypeOf(WebSocketTranscriptionsUpdatedEvent{}),
		string(EventTypeInputAudioBufferCompleted):      reflect.TypeOf(WebSocketInputAudioBufferCompletedEvent{}),
		string(EventTypeInputAudioBufferCleared):        reflect.TypeOf(WebSocketInputAudioBufferClearedEvent{}),
		string(EventTypeTranscriptionsMessageUpdate):    reflect.TypeOf(WebSocketTranscriptionsMessageUpdateEvent{}),
		string(EventTypeTranscriptionsMessageCompleted): reflect.TypeOf(WebSocketTranscriptionsMessageCompletedEvent{}),
		// v1/chat req
		string(EventTypeChatUpdate): reflect.TypeOf(WebSocketChatUpdateEvent{}),
		// string(EventTypeInputAudioBufferAppend):   reflect.TypeOf(WebSocketInputAudioBufferAppendEvent{}),
		// string(EventTypeInputAudioBufferComplete): reflect.TypeOf(WebSocketInputAudioBufferCompleteEvent{}),
		// string(EventTypeInputAudioBufferClear):    reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
		string(EventTypeConversationMessageCreate):         reflect.TypeOf(WebSocketConversationMessageCreateEvent{}),
		string(EventTypeConversationClear):                 reflect.TypeOf(WebSocketConversationClear{}),
		string(EventTypeConversationChatSubmitToolOutputs): reflect.TypeOf(WebSocketConversationChatSubmitToolOutputsEvent{}),
		string(EventTypeConversationChatCancel):            reflect.TypeOf(WebSocketConversationChatCancelEvent{}),
		// v1/chat resp
		string(EventTypeChatCreated):                    reflect.TypeOf(WebSocketChatCreatedEvent{}),
		string(EventTypeChatUpdated):                    reflect.TypeOf(WebSocketChatUpdatedEvent{}),
		string(EventTypeConversationChatCreated):        reflect.TypeOf(WebSocketConversationChatCreatedEvent{}),
		string(EventTypeConversationChatInProgress):     reflect.TypeOf(WebSocketConversationChatInProgressEvent{}),
		string(EventTypeConversationMessageDelta):       reflect.TypeOf(WebSocketConversationMessageDeltaEvent{}),
		string(EventTypeConversationAudioSentenceStart): reflect.TypeOf(WebSocketConversationAudioSentenceStartEvent{}),
		string(EventTypeConversationAudioDelta):         reflect.TypeOf(WebSocketConversationAudioDeltaEvent{}),
		string(EventTypeConversationMessageCompleted):   reflect.TypeOf(WebSocketConversationMessageCompletedEvent{}),
		string(EventTypeConversationAudioCompleted):     reflect.TypeOf(WebSocketConversationAudioCompletedEvent{}),
		string(EventTypeConversationChatCompleted):      reflect.TypeOf(WebSocketConversationChatCompletedEvent{}),
		string(EventTypeConversationChatFailed):         reflect.TypeOf(WebSocketConversationChatFailedEvent{}),
		// string(EventTypeInputAudioBufferCompleted):      reflect.TypeOf(WebSocketInputAudioBufferCompletedEvent{}),
		// string(EventTypeInputAudioBufferCleared):        reflect.TypeOf(WebSocketInputAudioBufferClearEvent{}),
		string(EventTypeConversationCleared):                  reflect.TypeOf(WebSocketConversationClearedEvent{}),
		string(EventTypeConversationChatCanceled):             reflect.TypeOf(WebSocketConversationChatCanceledEvent{}),
		string(EventTypeConversationAudioTranscriptUpdate):    reflect.TypeOf(WebSocketConversationAudioTranscriptUpdateEvent{}),
		string(EventTypeConversationAudioTranscriptCompleted): reflect.TypeOf(WebSocketConversationAudioTranscriptCompletedEvent{}),
		string(EventTypeConversationChatRequiresAction):       reflect.TypeOf(WebSocketConversationChatRequiresActionEvent{}),
		string(EventTypeInputAudioBufferSpeechStarted):        reflect.TypeOf(WebSocketInputAudioBufferSpeechStartedEvent{}),
		string(EventTypeInputAudioBufferSpeechStopped):        reflect.TypeOf(WebSocketInputAudioBufferSpeechStoppedEvent{}),
	}
}

// IWebSocketEvent websocket 事件接口
type IWebSocketEvent interface {
	GetEventType() WebSocketEventType
	GetID() string
	GetDetail() *EventDetail
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

	event := reflect.New(eventTypeRef).Interface().(IWebSocketEvent)
	if err := json.Unmarshal(common.Data, event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return event, nil
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
	OriginMessage string `json:"origin_message,omitempty"`
}
