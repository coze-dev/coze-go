package coze

import (
	"context"
	"net/http"
)

func (r *audioRooms) Create(ctx context.Context, req *CreateAudioRoomsReq) (*CreateAudioRoomsResp, error) {
	method := http.MethodPost
	uri := "/v1/audio/rooms"
	resp := &createAudioRoomsResp{}
	if err := r.core.Request(ctx, method, uri, req, resp); err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

type audioRooms struct {
	core *core
}

func newRooms(core *core) *audioRooms {
	return &audioRooms{core: core}
}

// AudioCodec represents the audio codec
type AudioCodec string

const (
	AudioCodecAACLC AudioCodec = "AACLC"
	AudioCodecG711A AudioCodec = "G711A"
	AudioCodecOPUS  AudioCodec = "OPUS"
	AudioCodecG722  AudioCodec = "G722"
)

// CreateAudioRoomsReq represents the request for creating an audio room
type CreateAudioRoomsReq struct {
	BotID          string      `json:"bot_id"`
	ConversationID string      `json:"conversation_id,omitempty"`
	VoiceID        string      `json:"voice_id,omitempty"`
	UID            string      `json:"uid,omitempty"`
	Config         *RoomConfig `json:"config,omitempty"`
}

// RoomConfig represents the room configuration
type RoomConfig struct {
	AudioConfig *RoomAudioConfig `json:"audio_config"`
}

// RoomAudioConfig represents the room audio configuration
type RoomAudioConfig struct {
	Codec AudioCodec `json:"codec"`
}

// createAudioRoomsResp represents the response for creating an audio room
type createAudioRoomsResp struct {
	baseResponse
	Data *CreateAudioRoomsResp `json:"data"`
}

// CreateAudioRoomsResp represents the response for creating an audio room
type CreateAudioRoomsResp struct {
	baseModel
	RoomID string `json:"room_id"`
	AppID  string `json:"app_id"`
	Token  string `json:"token"`
	UID    string `json:"uid"`
}
