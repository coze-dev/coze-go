package coze

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// IWebSocketEvent websocket 事件接口
type IWebSocketEvent interface {
	GetEventType() WebSocketEventType
	GetID() string
	GetDetail() *EventDetail
}

// WebSocketClientErrorEvent represents an client error event
// seq:common:1
type WebSocketClientErrorEvent struct {
	baseWebSocketEvent
	Data error `json:"data,omitempty"`
}

// WebSocketClosedEvent represents an closed event
// seq:common:2
type WebSocketClosedEvent struct {
	baseWebSocketEvent
}

// WebSocketErrorEvent represents an error event
// seq:common:3
type WebSocketErrorEvent struct {
	baseWebSocketEvent
	Data *Error `json:"data,omitempty"`
}

// v1/audio/speech req

// WebSocketSpeechUpdateEvent 流式输入文字
//
// 流式向服务端提交文字的片段。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#0ba93be3
// seq:v1/audio/speech:req:1
type WebSocketSpeechUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketSpeechUpdateEventData `json:"data,omitempty"`
}

// WebSocketSpeechUpdateEventData contains speech update configuration
type WebSocketSpeechUpdateEventData struct {
	// 输出音频格式。
	OutputAudio *OutputAudio `json:"output_audio,omitempty"`
}

// LimitConfig configures audio limits
type LimitConfig struct {
	// 周期的时长，单位为秒。例如设置为 10 秒，则以 10 秒作为一个周期。
	Period *int `json:"period,omitempty"`
	// 周期内，最大返回包数量。
	MaxFrameNum *int `json:"max_frame_num,omitempty"`
}

// InputAudio configuration for audio input
type InputAudio struct {
	// 输入音频的格式，支持 pcm、wav、ogg。默认为 wav。
	Format *string `json:"format,omitempty"`
	// 输入音频的编码，支持 pcm、opus、g711a、g711u。默认为 pcm。如果音频编码格式为 g711a 或 g711u，format 请设置为 pcm。
	Codec *string `json:"codec,omitempty"`
	// 输入音频的采样率，默认是 24000。支持 8000、16000、22050、24000、32000、44100、48000。如果音频编码格式 codec 为 g711a 或 g711u，音频采样率需设置为 8000。
	SampleRate *int `json:"sample_rate,omitempty"`
	// 输入音频的声道数，支持 1（单声道）、2（双声道）。默认是 1（单声道）。
	Channel *int `json:"channel,omitempty"`
	// 输入音频的位深，默认是 16，支持8、16和24。
	BitDepth *int `json:"bit_depth,omitempty"`
}

// OpusConfig configures Opus audio output
type OpusConfig struct {
	// 输出 opus 的码率，默认 48000。
	Bitrate *int `json:"bitrate,omitempty"`
	// 输出 opus 是否使用 CBR 编码，默认为 false。
	UseCBR *bool `json:"use_cbr,omitempty"`
	// 输出 opus 的帧长，默认是 10。可选值：2.5、5、10、20、40、60
	FrameSizeMs *float64 `json:"frame_size_ms,omitempty"`
	// 输出音频限流配置，默认不限制。
	LimitConfig *LimitConfig `json:"limit_config,omitempty"`
}

// PCMConfig configures PCM audio output
type PCMConfig struct {
	// 输出 pcm 音频的采样率，默认是 24000。支持 8000、16000、22050、24000、32000、44100、48000。
	SampleRate *int `json:"sample_rate,omitempty"`
	// 输出每个 pcm 包的时长，单位 ms，默认不限制。
	FrameSizeMs *float64 `json:"frame_size_ms,omitempty"`
	// 输出音频限流配置，默认不限制。
	LimitConfig *LimitConfig `json:"limit_config,omitempty"`
}

// OutputAudio configuration for audio output
type OutputAudio struct {
	// Output audio codec, supports pcm, g711a, g711u, opus. Default is pcm.
	Codec *string `json:"codec,omitempty"`
	// 当 codec 设置为 pcm、g711a 或 g711u 时，用于配置 PCM 音频参数。当 codec 设置为 opus 时，不需要设置此字段
	PCMConfig *PCMConfig `json:"pcm_config,omitempty"`
	// 当 codec 设置为 pcm 时，不需要设置此字段。
	OpusConfig *OpusConfig `json:"opus_config,omitempty"`
	// 输出音频的语速，取值范围 [-50, 100]，默认为 0。-50 表示 0.5 倍速，100 表示 2 倍速。
	SpeechRate *int `json:"speech_rate,omitempty"`
	// 输出音频的音色 ID，默认是柔美女友音色。你可以调用[查看音色列表](https://www.coze.cn/open/docs/developer_guides/list_voices) API 查看当前可用的所有音色 ID。
	VoiceID *string `json:"voice_id,omitempty"`
}

// WebSocketInputTextBufferAppendEvent 流式输入文字
//
// 流式向服务端提交文字的片段。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#0ba93be3
// seq:v1/audio/speech:req:2
type WebSocketInputTextBufferAppendEvent struct {
	baseWebSocketEvent
	Data *WebSocketInputTextBufferAppendEventData `json:"data,omitempty"`
}

// WebSocketInputTextBufferAppendEventData contains the text delta
type WebSocketInputTextBufferAppendEventData struct {
	// 需要合成语音的文字片段。
	Delta string `json:"delta"`
}

// WebSocketInputTextBufferCompleteEvent 提交文字
//
// 提交 append 的文本，发送后将收到 input_text_buffer.completed 的下行事件。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#ab24ada9
// seq:v1/audio/speech:req:3
type WebSocketInputTextBufferCompleteEvent struct {
	baseWebSocketEvent
}

type WebSocketInputTextBufferCompleteEventData struct{}

// WebSocketSpeechCreatedEvent 语音合成连接成功
//
// 语音合成连接成功后，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#23c0993e
// seq:v1/audio/speech:resp:1
type WebSocketSpeechCreatedEvent struct {
	baseWebSocketEvent
}

// WebSocketSpeechUpdatedEvent 配置更新完成
//
// 配置更新成功后，会返回最新的配置。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#a3a59fb4
// seq:v1/audio/speech:resp:2
type WebSocketSpeechUpdatedEvent struct {
	baseWebSocketEvent
	Data *WebSocketSpeechUpdatedEventData `json:"data,omitempty"`
}

// WebSocketSpeechUpdatedEventData contains speech session information
type WebSocketSpeechUpdatedEventData struct {
	// 输出音频格式。
	OutputAudio *OutputAudio `json:"output_audio,omitempty"`
}

// WebSocketInputTextBufferCompletedEvent input_text_buffer 提交完成
//
// 流式提交的文字完成后，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#cf5e0495
type WebSocketInputTextBufferCompletedEvent struct {
	baseWebSocketEvent
}

type WebSocketInputTextBufferCompletedEventData struct{}

// WebSocketSpeechAudioUpdateEvent 合成增量语音
//
// 语音合成产生增量语音时，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#98163c71
// seq:v1/audio/speech:resp:4
type WebSocketSpeechAudioUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketSpeechAudioUpdateEventData `json:"data,omitempty"`
}

// WebSocketSpeechAudioUpdateEventData contains audio delta
type WebSocketSpeechAudioUpdateEventData struct {
	// 音频片段。(API 返回的是base64编码的音频片段, SDK 已经自动解码为 bytes)
	Delta []byte `json:"delta"`
}

func (r WebSocketSpeechAudioUpdateEvent) dumpWithoutBinary() string {
	b, _ := json.Marshal(map[string]any{
		"type":   r.GetEventType(),
		"id":     r.GetID(),
		"detail": r.GetDetail(),
		"data": map[string]any{
			"delta": fmt.Sprintf("<length: %d>", len(r.Data.Delta)),
		},
	})
	return string(b)
}

// WebSocketSpeechAudioCompletedEvent 合成完成
//
// 语音合成完成后，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/tts_event#f42e9cb7
// seq:v1/audio/speech:resp:5
type WebSocketSpeechAudioCompletedEvent struct {
	baseWebSocketEvent
}

// v1/audio/transcriptions

// req

// WebSocketTranscriptionsUpdateEvent 更新语音识别配置
//
// 更新语音识别配置。若更新成功，会收到 transcriptions.updated 的下行事件，否则，会收到 error 下行事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#a7ca67ab
// seq:v1/audio/transcriptions:req:1
type WebSocketTranscriptionsUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketTranscriptionsUpdateEventData `json:"data,omitempty"`
}

// WebSocketTranscriptionsUpdateEventData contains transcription configuration
type WebSocketTranscriptionsUpdateEventData struct {
	// 输入音频格式。
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

// WebSocketInputAudioBufferAppendEvent 流式上传音频片段
//
// 流式向服务端提交音频的片段。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#9ef6e6ca
// seq:v1/audio/transcriptions:req:2
type WebSocketInputAudioBufferAppendEvent struct {
	baseWebSocketEvent
	Data *WebSocketInputAudioBufferAppendEventData `json:"data,omitempty"`
}

// WebSocketInputAudioBufferAppendEventData contains audio delta
type WebSocketInputAudioBufferAppendEventData struct {
	// 音频片段。(API 返回的是base64编码的音频片段, SDK 已经自动解码为 bytes)
	Delta []byte `json:"delta"`
}

// WebSocketInputAudioBufferCompleteEvent 提交音频
//
// 客户端发送 input_audio_buffer.complete 事件来告诉服务端提交音频缓冲区的数据。服务端提交成功后会返回 input_audio_buffer.completed 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#f5d76c87
// seq:v1/audio/transcriptions:req:3
type WebSocketInputAudioBufferCompleteEvent struct {
	baseWebSocketEvent
}

type WebSocketInputAudioBufferCompleteEventData struct{}

// WebSocketInputAudioBufferClearEvent 清除缓冲区音频
//
// 客户端发送 input_audio_buffer.clear 事件来告诉服务端清除缓冲区的音频数据。服务端清除完后将返回 input_audio_buffer.cleared 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#e98db543
// seq:v1/audio/transcriptions:req:4
type WebSocketInputAudioBufferClearEvent struct {
	baseWebSocketEvent
}

type WebSocketInputAudioBufferClearEventData struct{}

// WebSocketTranscriptionsCreatedEvent 语音识别连接成功
//
// 语音识别连接成功后，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#06d772a3
// seq:v1/audio/transcriptions:resp:1
type WebSocketTranscriptionsCreatedEvent struct {
	baseWebSocketEvent
}

// WebSocketTranscriptionsUpdatedEvent 配置更新成功
//
// 配置更新成功后，会返回最新的配置。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#3f842df1
// seq:v1/audio/transcriptions:resp:2
type WebSocketTranscriptionsUpdatedEvent struct {
	baseWebSocketEvent
	Data *WebSocketTranscriptionsUpdatedEventData `json:"data,omitempty"`
}

type WebSocketTranscriptionsUpdatedEventData struct {
	// 输入音频格式。
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

// WebSocketInputAudioBufferCompletedEvent 音频提交完成
//
// 客户端发送 input_audio_buffer.complete 事件来告诉服务端提交音频缓冲区的数据。服务端提交成功后会返回 input_audio_buffer.completed 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#8d747148
// seq:v1/audio/transcriptions:resp:3
type WebSocketInputAudioBufferCompletedEvent struct {
	baseWebSocketEvent
}

// WebSocketInputAudioBufferClearedEvent 音频清除成功
//
// 客户端发送 input_audio_buffer.clear 事件来告诉服务端清除音频缓冲区的数据。服务端清除完后将返回 input_audio_buffer.cleared 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#8211875b
// seq:v1/audio/transcriptions:resp:4
type WebSocketInputAudioBufferClearedEvent struct {
	baseWebSocketEvent
}

// WebSocketTranscriptionsMessageUpdateEvent 识别出文字
//
// 语音识别出文字后，返回此事件，每次都返回全量的识别出来的文字。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#772e6d2d
// seq:v1/audio/transcriptions:resp:5
type WebSocketTranscriptionsMessageUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketTranscriptionsMessageUpdateEventData `json:"data,omitempty"`
}

type WebSocketTranscriptionsMessageUpdateEventData struct {
	// 识别出的文字。
	Content string `json:"content"`
}

// WebSocketTranscriptionsMessageCompletedEvent 识别完成
//
// 语音识别完成后，返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/asr_event#0c36158c
type WebSocketTranscriptionsMessageCompletedEvent struct {
	baseWebSocketEvent
}

// v1/chat

// WebSocketChatUpdateEvent 更新对话配置
//
// 此事件可以更新当前对话连接的配置项，若更新成功，会收到 chat.updated 的下行事件，否则，会收到 error 下行事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#91642fa8
// seq:v1/chat:req:1
type WebSocketChatUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketChatUpdateEventData `json:"data,omitempty"`
}

// WebSocketChatUpdateEventData contains chat configuration
type WebSocketChatUpdateEventData struct {
	// 输出音频格式。
	OutputAudio *OutputAudio `json:"output_audio,omitempty"`
	// 输入音频格式。
	InputAudio *InputAudio `json:"input_audio,omitempty"`
	// 对话配置。
	ChatConfig *ChatConfig `json:"chat_config,omitempty"`
	// 需要订阅下行事件的事件类型列表。不设置或者设置为空为订阅所有下行事件。
	EventSubscriptions []string `json:"event_subscriptions,omitempty"`
	// 是否需要播放开场白，默认为 false。
	NeedPlayPrologue *bool `json:"need_play_prologue,omitempty"`
	// 自定义开场白，need_play_prologue 设置为 true 时生效。如果不设定自定义开场白则使用智能体上设置的开场白。
	PrologueContent string `json:"prologue_content,omitempty"`
	// 转检测配置。
	TurnDetection *TurnDetection `json:"turn_detection,omitempty"`
	// 语音识别配置，包括热词和上下文信息，以便优化语音识别的准确性和相关性。
	AsrConfig *ASRConfig `json:"asr_config,omitempty"`
}

type ChatConfig struct {
	// 标识对话发生在哪一次会话中。会话是智能体和用户之间的一段问答交互。一个会话包含一条或多条消息。对话是会话中对智能体的一次调用，智能体会将对话中产生的消息添加到会话中。可以使用已创建的会话，会话中已存在的消息将作为上下文传递给模型。创建会话的方式可参考创建会话。对于一问一答等不需要区分 conversation 的场合可不传该参数，系统会自动生成一个会话。不传的话会默认创建一个新的 conversation。
	ConversationID *string `json:"conversation_id,omitempty"`
	// 标识当前与智能体的用户，由使用方自行定义、生成与维护。user_id 用于标识对话中的不同用户，不同的 user_id，其对话的上下文消息、数据库等对话记忆数据互相隔离。如果不需要用户数据隔离，可将此参数固定为一个任意字符串，例如 123，abc 等。
	UserID *string `json:"user_id,omitempty"`
	// 附加信息，通常用于封装一些业务相关的字段。查看对话消息详情时，系统会透传此附加信息。自定义键值对，应指定为 Map 对象格式。长度为 16 对键值对，其中键（key）的长度范围为 1～64 个字符，值（value）的长度范围为 1～512 个字符。
	MetaData map[string]string `json:"meta_data,omitempty"`
	// 智能体中定义的变量。在智能体 prompt 中设置变量 {{key}} 后，可以通过该参数传入变量值，同时支持 Jinja2 语法。详细说明可参考变量示例。变量名只支持英文字母和下划线。
	CustomVariables map[string]string `json:"custom_variables,omitempty"`
	// 附加参数，通常用于特殊场景下指定一些必要参数供模型判断，例如指定经纬度，并询问智能体此位置的天气。自定义键值对格式，其中键（key）仅支持设置为：latitude（纬度，此时值（Value）为纬度值，例如 39.9800718）。longitude（经度，此时值（Value）为经度值，例如 116.309314）。
	ExtraParams map[string]string `json:"extra_params,omitempty"`
	// 是否保存本次对话记录。true：（默认）会话中保存本次对话记录，包括本次对话的模型回复结果、模型执行中间结果。false：会话中不保存本次对话记录，后续也无法通过任何方式查看本次对话信息、消息详情。在同一个会话中再次发起对话时，本次会话也不会作为上下文传递给模型。
	AutoSaveHistory *bool `json:"auto_save_history,omitempty"`
	// 设置对话流的自定义输入参数的值，具体用法和示例代码可参考[为自定义参数赋值](https://www.coze.cn/open/docs/tutorial/variable)。 对话流的输入参数 USER_INPUT 应在 additional_messages 中传入，在 parameters 中的 USER_INPUT 不生效。 如果 parameters 中未指定 CONVERSATION_NAME 或其他输入参数，则使用参数默认值运行对话流；如果指定了这些参数，则使用指定值。
	Parameters map[string]any `json:"parameters,omitempty"`
}

type TurnDetection struct {
	// 用户演讲检测模式
	Type *TurnDetectionType `json:"type,omitempty"`
	// server_vad 模式下，VAD 检测到语音之前要包含的音频量，单位为 ms。默认为 600ms。
	PrefixPaddingMS *int64 `json:"prefix_padding_ms,omitempty"`
	// server_vad 模式下，检测语音停止的静音持续时间，单位为 ms。默认为 500ms。
	SilenceDurationMS *int64 `json:"silence_duration_ms,omitempty"`
	// server_vad 模式下打断策略配置
	InterruptConfig *InterruptConfig `json:"interrupt_config,omitempty"`
}

type TurnDetectionType string

const (
	// TurnDetectionTypeServerVAD 自由对话模式，语音数据会传输到服务器端进行实时分析，服务器端的语音活动检测算法会判断用户是否在说话。
	TurnDetectionTypeServerVAD TurnDetectionType = "server_vad"
	// TurnDetectionTypeClientInterrupt 按键说话模式，客户端实时分析语音数据，并检测用户是否已停止说话。
	TurnDetectionTypeClientInterrupt TurnDetectionType = "client_interrupt"
)

type InterruptConfig struct {
	// 打断模式
	Mode *InterruptConfigMode `json:"mode,omitempty"`
	// 打断的关键词配置，最多同时限制 5 个关键词，每个关键词限定长度在6-24个字节以内(2-8个汉字以内), 不能有标点符号。
	Keywords []string `json:"keywords,omitempty"`
}

type InterruptConfigMode string

const (
	// InterruptConfigModeKeywordContains keyword_contains模式下，说话内容包含关键词才会打断模型回复。例如关键词"扣子"，用户正在说“你好呀扣子......” / “扣子你好呀”，模型回复都会被打断。
	InterruptConfigModeKeywordContains InterruptConfigMode = "keyword_contains"
	// InterruptConfigModeKeywordPrefix keyword_prefix模式下，说话内容前缀匹配关键词才会打断模型回复。例如关键词"扣子"，用户正在说“扣子你好呀......”，模型回复就会被打断，而用户说“你好呀扣子......”，模型回复不会被打断。
	InterruptConfigModeKeywordPrefix InterruptConfigMode = "keyword_prefix"
)

type ASRConfig struct {
	// 请输入热词列表，以便提升这些词汇的识别准确率。所有热词加起来最多100个 Tokens，超出部分将自动截断。
	HotWords []string `json:"hot_words,omitempty"`
	// 请输入上下文信息。最多输入 800 个 Tokens，超出部分将自动截断。
	Context *string `json:"context,omitempty"`
	// 用户说话的语种，默认为 common。选项包括：
	UserLanguage *ASRConfigUserLanguage `json:"user_language,omitempty"`
	// 将语音转为文本时，是否启用语义顺滑。默认为 true。true：系统在进行语音处理时，会去掉识别结果中诸如 “啊”“嗯” 等语气词，使得输出的文本语义更加流畅自然，符合正常的语言表达习惯，尤其适用于对文本质量要求较高的场景，如正式的会议记录、新闻稿件生成等。false：系统不会对识别结果中的语气词进行处理，识别结果会保留原始的语气词。
	EnableDDC *bool `json:"enable_ddc,omitempty"`
	// 将语音转为文本时，是否开启文本规范化（ITN）处理，将识别结果转换为更符合书面表达习惯的格式以提升可读性。默认为 true。开启后，会将口语化数字转换为标准数字格式，示例：将两点十五分转换为 14:15。将一百美元转换为 $100。
	EnableITN *bool `json:"enable_itn,omitempty"`
	// 将语音转为文本时，是否给文本加上标点符号。默认为 true。
	EnablePunc *bool `json:"enable_punc,omitempty"`
}

type ASRConfigUserLanguage string

const (
	ASRConfigUserLanguageCommon ASRConfigUserLanguage = "common" // 大模型语音识别，可自动识别中英粤。
	ASRConfigUserLanguageZH     ASRConfigUserLanguage = "zh"     // 小模型语音识别，中文。
	ASRConfigUserLanguageCANT   ASRConfigUserLanguage = "cant"   // 小模型语音识别，粤语。
	ASRConfigUserLanguageSC     ASRConfigUserLanguage = "sc"     // 小模型语音识别，川渝。
	ASRConfigUserLanguageEN     ASRConfigUserLanguage = "en"     // 小模型语音识别，英语。
	ASRConfigUserLanguageJA     ASRConfigUserLanguage = "ja"     // 小模型语音识别，日语。
	ASRConfigUserLanguageKO     ASRConfigUserLanguage = "ko"     // 小模型语音识别，韩语。
	ASRConfigUserLanguageFR     ASRConfigUserLanguage = "fr"     // 小模型语音识别，法语。
	ASRConfigUserLanguageID     ASRConfigUserLanguage = "id"     // 小模型语音识别，印尼语。
	ASRConfigUserLanguageES     ASRConfigUserLanguage = "es"     // 小模型语音识别，西班牙语。
	ASRConfigUserLanguagePT     ASRConfigUserLanguage = "pt"     // 小模型语音识别，葡萄牙语。
	ASRConfigUserLanguageMS     ASRConfigUserLanguage = "ms"     // 小模型语音识别，马来语。
	ASRConfigUserLanguageRU     ASRConfigUserLanguage = "ru"     // 小模型语音识别，俄语。
)

// WebSocketConversationMessageCreateEvent 手动提交对话内容
//
// 若 role=user，提交事件后就会生成语音回复，适合如下的场景，比如帮我解析 xx 链接，帮我分析这个图片的内容等。若 role=assistant，提交事件后会加入到对话的上下文。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#46f6a7d0
// seq:v1/chat:req:2
type WebSocketConversationMessageCreateEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationMessageCreateEventData `json:"data,omitempty"`
}

// WebSocketConversationMessageCreateEventData contains message content
type WebSocketConversationMessageCreateEventData struct {
	// 发送这条消息的实体。取值：user（代表该条消息内容是用户发送的）、assistant（代表该条消息内容是智能体发送的）。
	Role MessageRole `json:"role,omitempty"`
	// 消息内容的类型，支持设置为：text：文本。object_string：多模态内容，即文本和文件的组合、文本和图片的组合。
	ContentType MessageContentType `json:"content_type,omitempty"`
	// 消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）、卡片等多种类型的内容。当 content_type 为 object_string时，content 的结构和详细参数说明请参见object_string object。
	Content string `json:"content,omitempty"`
}

// WebSocketConversationClearEvent 清除上下文
//
// 清除上下文，会在当前 conversation 下新增一个 section，服务端处理完后会返回 conversation.cleared 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#aa86f213
type WebSocketConversationClearEvent struct {
	baseWebSocketEvent
}

type WebSocketConversationClearEventData struct{}

// WebSocketConversationChatSubmitToolOutputsEvent 提交端插件执行结果
//
// 你可以将需要客户端执行的操作定义为插件，对话中如果触发这个插件，会收到一个 event_type = "conversation.chat.requires_action" 的下行事件，此时需要执行客户端的操作后，通过此上行事件来提交插件执行后的结果。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#aacdcb41
// seq:v1/chat:req:3
type WebSocketConversationChatSubmitToolOutputsEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationChatSubmitToolOutputsEventData `json:"data,omitempty"`
}

// WebSocketConversationChatSubmitToolOutputsEventData contains tool outputs
type WebSocketConversationChatSubmitToolOutputsEventData struct {
	// 对话的唯一标识。
	ChatID string `json:"chat_id"`
	// 工具执行结果。
	ToolOutputs []*ToolOutput `json:"tool_outputs"`
}

// WebSocketConversationChatCancelEvent 打断智能体输出
//
// 发送此事件可取消正在进行的对话，中断后，服务端将会返回 conversation.chat.canceled 事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#0554db7d
// seq:v1/chat:req:4
type WebSocketConversationChatCancelEvent struct {
	baseWebSocketEvent
}

type WebSocketConversationChatCancelEventData struct{}

// WebSocketInputTextGenerateAudioEvent 文本生成语音
//
// 你可以主动提交一段文字用来做语音合成，提交的消息不会触发智能体的回复，只会合成音频内容下发到客户端。
// 提交事件的时候如果智能体正在输出语音会被中断输出。
// 适合在和智能体聊天过程中客户端长时间没有响应，智能体可以主动说话暖场的场景。
type WebSocketInputTextGenerateAudioEvent struct {
	baseWebSocketEvent
	Data *WebSocketInputTextGenerateAudioEventData `json:"data,omitempty"`
}

// WebSocketInputTextGenerateAudioEventData contains text to audio data
type WebSocketInputTextGenerateAudioEventData struct {
	// 消息内容的类型，支持设置为：text：文本
	Mode WebSocketInputTextGenerateAudioEventDataMode `json:"mode,omitempty"`
	// 当 mode == text 时候必填。长度限制 (0, 1024) 字节
	Text string `json:"text,omitempty"`
}

type WebSocketInputTextGenerateAudioEventDataMode string

const (
	WebSocketInputTextGenerateAudioEventDataModeText WebSocketInputTextGenerateAudioEventDataMode = "text"
)

// WebSocketChatCreatedEvent 对话连接成功
//
// 流式对话接口成功建立连接后服务端会发送此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#a061f115
// seq:v1/chat:resp:1
type WebSocketChatCreatedEvent struct {
	baseWebSocketEvent
}

// WebSocketChatUpdatedEvent 对话配置成功
//
// 对话配置更新成功后，会返回最新的配置。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#39879618
// seq:v1/chat:resp:2
type WebSocketChatUpdatedEvent struct {
	baseWebSocketEvent
	Data *WebSocketChatUpdateEventData `json:"data,omitempty"`
}

// WebSocketConversationChatCreatedEvent 对话开始
//
// 创建对话的事件，表示对话开始。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#a2b10fd2
// seq:v1/chat:resp:3
type WebSocketConversationChatCreatedEvent struct {
	baseWebSocketEvent
	Data *Chat `json:"data,omitempty"`
}

// WebSocketConversationChatInProgressEvent 对话正在处理
//
// 服务端正在处理对话。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#36a38a6b
// seq:v1/chat:resp:4
type WebSocketConversationChatInProgressEvent struct {
	baseWebSocketEvent
	Data *Chat `json:"data,omitempty"`
}

// WebSocketConversationMessageDeltaEvent 增量消息
//
// 增量消息，通常是 type=answer 时的增量消息。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#2dfe8dba
// seq:v1/chat:resp:5
type WebSocketConversationMessageDeltaEvent struct {
	baseWebSocketEvent
	Data *Message `json:"data,omitempty"`
}

// WebSocketConversationAudioSentenceStartEvent 增量语音字幕
//
// 一条新的字幕句子，后续的 conversation.audio.delta 增量语音均属于当前字幕句子，可能有多个增量语音共同对应此句字幕的文字内容。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#2e67bf44
// seq:v1/chat:resp:6
type WebSocketConversationAudioSentenceStartEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationAudioSentenceStartEventData `json:"data,omitempty"`
}

type WebSocketConversationAudioSentenceStartEventData struct {
	// 新字幕句子的文本内容，后续相关 conversation.audio.delta 增量语音均对应此文本。
	Text string `json:"text"`
}

// WebSocketConversationAudioDeltaEvent 增量语音
//
// 增量消息，通常是 type=answer 时的增量消息。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#36a38a6b
// seq:v1/chat:resp:7
type WebSocketConversationAudioDeltaEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationAudioDeltaEventData `json:"data,omitempty"`
}

type WebSocketConversationAudioDeltaEventData struct {
	// The entity that sent this message.
	Role MessageRole `json:"role"`

	// The type of message.
	Type MessageType `json:"type"`

	// The content of the message. It supports various types of content, including plain text,
	// multimodal (a mix of text, images, and files), message cards, and more.
	Content []byte `json:"content"`

	// The reasoning_content of the thought process message
	ReasoningContent string `json:"reasoning_content"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages. Custom key-value pairs should be specified in Map object
	// format, with a length of 16 key-value pairs. The length of the key should be between 1 and 64
	// characters, and the length of the value should be between 1 and 512 characters.
	MetaData map[string]string `json:"meta_data,omitempty"`

	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`

	// section_id is used to distinguish the context sections of the session history. The same section
	// is one context.
	SectionID string `json:"section_id"`
	BotID     string `json:"bot_id"`
	ChatID    string `json:"chat_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (r *WebSocketConversationAudioDeltaEvent) dumpWithoutBinary() string {
	if r.Data.ContentType == MessageContentTypeAudio {
		b, _ := json.Marshal(map[string]any{
			"id":         r.GetID(),
			"event_type": r.GetEventType(),
			"detail":     r.GetDetail(),
			"data": &Message{
				Role:             r.Data.Role,
				Type:             r.Data.Type,
				Content:          fmt.Sprintf("<length: %d>", len(r.Data.Content)),
				ReasoningContent: r.Data.ReasoningContent,
				ContentType:      r.Data.ContentType,
				MetaData:         r.Data.MetaData,
				ID:               r.Data.ID,
				ConversationID:   r.Data.ConversationID,
				SectionID:        r.Data.SectionID,
				BotID:            r.Data.BotID,
				ChatID:           r.Data.ChatID,
				CreatedAt:        r.Data.CreatedAt,
				UpdatedAt:        r.Data.UpdatedAt,
			},
		})
		return string(b)
	}
	b, _ := json.Marshal(r)
	return string(b)
}

// WebSocketConversationMessageCompletedEvent 消息完成
//
// 消息已回复完成。此时事件中带有所有 message.delta 的拼接结果，且每个消息均为 completed 状态。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#4361e8d1
// seq:v1/chat:resp:8
type WebSocketConversationMessageCompletedEvent struct {
	baseWebSocketEvent
	Data *Message `json:"data,omitempty"`
}

// WebSocketConversationAudioCompletedEvent 语音回复完成
//
// 语音回复完成，表示对话结束。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#b00d6a73
// seq:v1/chat:resp:9
type WebSocketConversationAudioCompletedEvent struct {
	baseWebSocketEvent
	Data *Message `json:"data,omitempty"`
}

// WebSocketConversationChatCompletedEvent 对话完成
//
// 对话完成，表示对话结束。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#02fac327
// seq:v1/chat:resp:10
type WebSocketConversationChatCompletedEvent struct {
	baseWebSocketEvent
	Data *Chat `json:"data,omitempty"`
}

// WebSocketConversationChatFailedEvent 对话失败
//
// 此事件用于标识对话失败。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#765bb7e5
// seq:v1/chat:resp:11
type WebSocketConversationChatFailedEvent struct {
	baseWebSocketEvent
	Data *Chat `json:"data,omitempty"`
}

// WebSocketConversationClearedEvent 上下文清除完成
//
// 上下文清除完成，表示上下文已清除。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#6a941b8a
// seq:v1/chat:resp:12
type WebSocketConversationClearedEvent struct {
	baseWebSocketEvent
}

// WebSocketConversationChatCanceledEvent 智能体输出中断
//
// 客户端提交 conversation.chat.cancel 事件，服务端完成中断后，将返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#089ed144
// seq:v1/chat:resp:13
type WebSocketConversationChatCanceledEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationChatCanceledEventData `json:"data,omitempty"`
}

// WebSocketConversationChatCanceledEventData contains cancellation information
type WebSocketConversationChatCanceledEventData struct {
	// 输出中断类型枚举值，包括 1: 被用户语音说话打断  2: 用户主动 cancel  3: 手动提交对话内容
	Code int `json:"code"`
	// 智能体输出中断的详细说明
	Msg string `json:"msg"`
}

// WebSocketConversationAudioTranscriptUpdateEvent 用户语音识别字幕
//
// 用户语音识别的中间值，每次返回都是全量文本。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#1b59cbf9
// seq:v1/chat:resp:14
type WebSocketConversationAudioTranscriptUpdateEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationAudioTranscriptUpdateEventData `json:"data,omitempty"`
}

type WebSocketConversationAudioTranscriptUpdateEventData struct {
	// 语音识别的中间值。
	Content string `json:"content"`
}

// WebSocketConversationAudioTranscriptCompletedEvent 用户语音识别完成
//
// 用户语音识别完成，表示用户语音识别完成。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#9d1e6930
// seq:v1/chat:resp:15
type WebSocketConversationAudioTranscriptCompletedEvent struct {
	baseWebSocketEvent
	Data *WebSocketConversationAudioTranscriptCompletedEventData `json:"data,omitempty"`
}

type WebSocketConversationAudioTranscriptCompletedEventData struct {
	// 语音识别的最终结果。
	Content string `json:"content"`
}

// WebSocketConversationChatRequiresActionEvent 端插件请求
//
// 对话中断，需要使用方上报工具的执行结果。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#2ef697d8
// seq:v1/chat:resp:16
type WebSocketConversationChatRequiresActionEvent struct {
	baseWebSocketEvent
	Data *Chat `json:"data,omitempty"`
}

// WebSocketInputAudioBufferSpeechStartedEvent 用户开始说话
//
// 此事件表示服务端识别到用户正在说话。只有在 server_vad 模式下，才会返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#95553c68
// seq:v1/chat:resp:17
type WebSocketInputAudioBufferSpeechStartedEvent struct {
	baseWebSocketEvent
}

// WebSocketInputAudioBufferSpeechStoppedEvent 用户结束说话
//
// 此事件表示服务端识别到用户已停止说话。只有在 server_vad 模式下，才会返回此事件。
// docs: https://www.coze.cn/open/docs/developer_guides/streaming_chat_event#5084c0aa
// seq:v1/chat:resp:18
type WebSocketInputAudioBufferSpeechStoppedEvent struct {
	baseWebSocketEvent
}

type IWebSocketChatHandler interface {
	OnClientError(ctx context.Context, cli *WebSocketChat, event *WebSocketClientErrorEvent) error
	OnClosed(ctx context.Context, cli *WebSocketChat, event *WebSocketClosedEvent) error
	OnError(ctx context.Context, cli *WebSocketChat, event *WebSocketErrorEvent) error
	OnChatCreated(ctx context.Context, cli *WebSocketChat, event *WebSocketChatCreatedEvent) error
	OnChatUpdated(ctx context.Context, cli *WebSocketChat, event *WebSocketChatUpdatedEvent) error
	OnConversationChatCreated(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCreatedEvent) error
	OnConversationChatInProgress(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatInProgressEvent) error
	OnConversationMessageDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageDeltaEvent) error
	OnConversationAudioSentenceStart(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioSentenceStartEvent) error
	OnConversationAudioDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioDeltaEvent) error
	OnConversationMessageCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageCompletedEvent) error
	OnConversationAudioCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioCompletedEvent) error
	OnConversationChatCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCompletedEvent) error
	OnConversationChatFailed(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatFailedEvent) error
	OnInputAudioBufferCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferCompletedEvent) error
	OnInputAudioBufferCleared(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferClearedEvent) error
	OnConversationCleared(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationClearedEvent) error
	OnConversationChatCanceled(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCanceledEvent) error
	OnConversationAudioTranscriptUpdate(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptUpdateEvent) error
	OnConversationAudioTranscriptCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptCompletedEvent) error
	OnConversationChatRequiresAction(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatRequiresActionEvent) error
	OnInputAudioBufferSpeechStarted(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStartedEvent) error
	OnInputAudioBufferSpeechStopped(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStoppedEvent) error
}

type BaseWebSocketChatHandler struct{}

func (BaseWebSocketChatHandler) OnClientError(ctx context.Context, cli *WebSocketChat, event *WebSocketClientErrorEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnClosed(ctx context.Context, cli *WebSocketChat, event *WebSocketClosedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnError(ctx context.Context, cli *WebSocketChat, event *WebSocketErrorEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnChatCreated(ctx context.Context, cli *WebSocketChat, event *WebSocketChatCreatedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnChatUpdated(ctx context.Context, cli *WebSocketChat, event *WebSocketChatUpdatedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatCreated(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCreatedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatInProgress(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatInProgressEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationMessageDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageDeltaEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationAudioSentenceStart(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioSentenceStartEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationAudioDelta(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioDeltaEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationMessageCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationMessageCompletedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationAudioCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioCompletedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCompletedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatFailed(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatFailedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnInputAudioBufferCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferCompletedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnInputAudioBufferCleared(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferClearedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationCleared(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationClearedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatCanceled(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatCanceledEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationAudioTranscriptUpdate(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptUpdateEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationAudioTranscriptCompleted(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationAudioTranscriptCompletedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnConversationChatRequiresAction(ctx context.Context, cli *WebSocketChat, event *WebSocketConversationChatRequiresActionEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnInputAudioBufferSpeechStarted(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStartedEvent) error {
	return nil
}

func (BaseWebSocketChatHandler) OnInputAudioBufferSpeechStopped(ctx context.Context, cli *WebSocketChat, event *WebSocketInputAudioBufferSpeechStoppedEvent) error {
	return nil
}

type IWebSocketAudioSpeechHandler interface {
	OnClientError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error
	OnClosed(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error
	OnError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error
	OnSpeechCreated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error
	OnSpeechUpdated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error
	OnInputTextBufferCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error
	OnSpeechAudioUpdate(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error
	OnSpeechAudioCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error
}

type BaseWebSocketAudioSpeechHandler struct{}

func (BaseWebSocketAudioSpeechHandler) OnClientError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClientErrorEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnClosed(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketClosedEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnError(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketErrorEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnSpeechCreated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechCreatedEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnSpeechUpdated(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechUpdatedEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnInputTextBufferCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketInputTextBufferCompletedEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnSpeechAudioUpdate(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioUpdateEvent) error {
	return nil
}

func (BaseWebSocketAudioSpeechHandler) OnSpeechAudioCompleted(ctx context.Context, cli *WebSocketAudioSpeech, event *WebSocketSpeechAudioCompletedEvent) error {
	return nil
}

type IWebSocketAudioTranscriptionHandler interface {
	OnClientError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error
	OnClosed(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error
	OnError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error
	OnTranscriptionsCreated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error
	OnTranscriptionsUpdated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error
	OnInputAudioBufferCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error
	OnInputAudioBufferCleared(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error
	OnTranscriptionsMessageUpdate(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error
	OnTranscriptionsMessageCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error
}

type BaseWebSocketAudioTranscriptionHandler struct{}

func (BaseWebSocketAudioTranscriptionHandler) OnClientError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClientErrorEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnClosed(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketClosedEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnError(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketErrorEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnTranscriptionsCreated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsCreatedEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnTranscriptionsUpdated(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsUpdatedEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnInputAudioBufferCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferCompletedEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnInputAudioBufferCleared(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketInputAudioBufferClearedEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnTranscriptionsMessageUpdate(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageUpdateEvent) error {
	return nil
}

func (BaseWebSocketAudioTranscriptionHandler) OnTranscriptionsMessageCompleted(ctx context.Context, cli *WebSocketAudioTranscription, event *WebSocketTranscriptionsMessageCompletedEvent) error {
	return nil
}

func newWebSocketEvent(eventType WebSocketEventType, data any) IWebSocketEvent {
	eventStructType, exists := websocketEvents[string(eventType)]
	if !exists {
		return nil
	}
	eventValue := reflect.New(eventStructType).Elem()

	eventTypeField := eventValue.FieldByName("EventType")
	if eventTypeField.IsValid() && eventTypeField.CanSet() {
		eventTypeField.SetString(string(eventType))
	}
	dataField := eventValue.FieldByName("Data")
	if dataField.IsValid() && dataField.CanSet() {
		dataField.Set(reflect.ValueOf(data))
	}
	return eventValue.Interface().(IWebSocketEvent)
}
