package coze

import (
	"context"
	"net/http"
)

func (r *audioRooms) Create(ctx context.Context, req *CreateAudioRoomsReq) (*CreateAudioRoomsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/rooms",
		Body:   req,
	}
	response := new(createAudioRoomsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// AudioCodec represents the audio codec
type AudioCodec string

const (
	AudioCodecAACLC AudioCodec = "AACLC"
	AudioCodecG711A AudioCodec = "G711A"
	AudioCodecOPUS  AudioCodec = "OPUS"
	AudioCodecG722  AudioCodec = "G722"
)

func (r AudioCodec) String() string {
	return string(r)
}

func (r AudioCodec) Ptr() *AudioCodec {
	return &r
}

// CreateAudioRoomsReq represents the request for creating an audio room
type CreateAudioRoomsReq struct {
	BotID          string      `json:"bot_id"`
	ConversationID string      `json:"conversation_id,omitempty"`
	VoiceID        string      `json:"voice_id,omitempty"`
	UID            string      `json:"uid,omitempty"`
	WorkflowID     string      `json:"workflow_id,omitempty"`
	Config         *RoomConfig `json:"config,omitempty"`
}

// RoomConfig represents the room configuration
type RoomConfig struct {
	AudioConfig     *RoomAudioConfig `json:"audio_config,omitempty"`
	VideoConfig     *RoomVideoConfig `json:"video_config,omitempty"`
	PrologueContent string           `json:"prologue_content,omitempty"`
}

// VideoCodec represents the video codec
type VideoCodec string

const (
	VideoCodecH264    VideoCodec = "H264"
	VideoCodecBYTEVC1 VideoCodec = "BYTEVC1"
)

func (r VideoCodec) String() string {
	return string(r)
}

func (r VideoCodec) Ptr() *VideoCodec {
	return &r
}

// StreamVideoType represents the stream video type
type StreamVideoType string

const (
	// StreamVideoTypeMain 主流，包括通过摄像头/麦克风的内部采集机制获取的流，以及通过自定义采集方式获取的流。
	StreamVideoTypeMain StreamVideoType = "main"
	// StreamVideoTypeScreen 屏幕流，用于屏幕共享或屏幕录制的视频流。
	StreamVideoTypeScreen StreamVideoType = "screen"
)

func (r StreamVideoType) String() string {
	return string(r)
}

func (r StreamVideoType) Ptr() *StreamVideoType {
	return &r
}

// RoomVideoConfig represents the room video configuration
type RoomVideoConfig struct {
	// 房间视频编码格式
	Codec VideoCodec `json:"codec,omitempty"`
	// 房间视频流类型
	StreamVideoType StreamVideoType `json:"stream_video_type,omitempty"`
	// 视频抽帧速率
	VideoFrameRate *int `json:"video_frame_rate,omitempty"`
	// 视频帧过期时间, 单位为 s
	VideoFrameExpireDuration *int `json:"video_frame_expire_duration,omitempty"`
}

// RoomAudioConfig represents the room audio configuration
type RoomAudioConfig struct {
	Codec AudioCodec `json:"codec"`
}

// CreateAudioRoomsResp represents the response for creating an audio room
type CreateAudioRoomsResp struct {
	baseModel
	RoomID string `json:"room_id"`
	AppID  string `json:"app_id"`
	Token  string `json:"token"`
	UID    string `json:"uid"`
}

type createAudioRoomsResp struct {
	baseResponse
	Data *CreateAudioRoomsResp `json:"data"`
}

type audioRooms struct {
	core *core
}

func newRooms(core *core) *audioRooms {
	return &audioRooms{core: core}
}
