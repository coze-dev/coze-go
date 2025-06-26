package coze

import (
	"context"
	"io"
	"net/http"
)

func (r *audioVoices) Clone(ctx context.Context, req *CloneAudioVoicesReq) (*CloneAudioVoicesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/voices/clone",
		Body:   req,
		IsFile: true,
	}
	response := new(cloneAudioVoicesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *audioVoices) List(ctx context.Context, req *ListAudioVoicesReq) (NumberPaged[Voice], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[Voice], error) {
			response := &ListAudioVoicesResp{}
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/audio/voices",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[Voice]{
				response: response.HTTPResponse,
				HasMore:  len(response.Data.VoiceList) >= request.PageSize,
				Data:     response.Data.VoiceList,
				LogID:    response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type VoiceState string

const (
	VoiceStateInit   VoiceState = "init"   // 初始化
	VoiceStateCloned VoiceState = "cloned" // 已克隆
	VoiceStateAll    VoiceState = "all"    // 所有, 只有查询的时候有效
)

func (r VoiceState) String() string {
	return string(r)
}

func (r VoiceState) Ptr() *VoiceState {
	return &r
}

type VoiceModelType string

const (
	VoiceModelTypeBig   VoiceModelType = "big"   // 大模型音色
	VoiceModelTypeSmall VoiceModelType = "small" // 小模型音色
)

func (r VoiceModelType) String() string {
	return string(r)
}

func (r VoiceModelType) Ptr() *VoiceModelType {
	return &r
}

// Voice represents the voice model
type Voice struct {
	VoiceID                string         `json:"voice_id"`
	Name                   string         `json:"name"`
	IsSystemVoice          bool           `json:"is_system_voice"`
	LanguageCode           string         `json:"language_code"`
	LanguageName           string         `json:"language_name"`
	PreviewText            string         `json:"preview_text"`
	PreviewAudio           string         `json:"preview_audio"`
	AvailableTrainingTimes int            `json:"available_training_times"`
	CreateTime             int            `json:"create_time"`
	UpdateTime             int            `json:"update_time"`
	ModelType              VoiceModelType `json:"model_type"`
	State                  VoiceState     `json:"state"`
}

// CloneAudioVoicesReq represents the request for cloning a voice
type CloneAudioVoicesReq struct {
	VoiceName   string        `json:"voice_name"`
	File        io.Reader     `json:"file"`
	AudioFormat AudioFormat   `json:"audio_format"`
	Language    *LanguageCode `json:"language"`
	VoiceID     *string       `json:"voice_id"`
	PreviewText *string       `json:"preview_text"`
	Text        *string       `json:"text"`
	SpaceID     *string       `json:"space_id"`
	Description *string       `json:"description"`
}

// CloneAudioVoicesResp represents the response for cloning a voice
type CloneAudioVoicesResp struct {
	baseModel
	VoiceID string `json:"voice_id"`
}

// ListAudioVoicesReq represents the request for listing voices
type ListAudioVoicesReq struct {
	FilterSystemVoice bool            `query:"filter_system_voice" json:"-"`
	PageNum           int             `query:"page_num" json:"-"`
	PageSize          int             `query:"page_size" json:"-"`
	ModelType         *VoiceModelType `query:"model_type" json:"-"`
	VoiceState        *VoiceState     `query:"voice_state" json:"-"`
}

// ListAudioVoicesResp represents the response for listing voices
type ListAudioVoicesResp struct {
	baseResponse
	Data struct {
		VoiceList []*Voice `json:"voice_list"`
	} `json:"data"`
}

type cloneAudioVoicesResp struct {
	baseResponse
	Data *CloneAudioVoicesResp `json:"data"`
}

func (r ListAudioVoicesReq) toReq(request *pageRequest) *ListAudioVoicesReq {
	return &ListAudioVoicesReq{
		FilterSystemVoice: r.FilterSystemVoice,
		ModelType:         r.ModelType,
		VoiceState:        r.VoiceState,
		PageNum:           request.PageNum,
		PageSize:          request.PageSize,
	}
}

type audioVoices struct {
	core *core
}

func newAudioVoices(core *core) *audioVoices {
	return &audioVoices{core: core}
}
