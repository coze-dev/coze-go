package coze

import (
	"context"
	"io"
	"net/http"
)

func (r *audioTranscriptions) Create(ctx context.Context, req *AudioSpeechTranscriptionsReq) (*CreateAudioTranscriptionsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/transcriptions",
		Body:   req,
		IsFile: true,
	}
	response := new(createAudioTranscriptionsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.CreateAudioTranscriptionsResp, err
}

type AudioSpeechTranscriptionsReq struct {
	Filename string    `json:"filename"`
	Audio    io.Reader `json:"file"`
}

type createAudioTranscriptionsResp struct {
	baseResponse
	*CreateAudioTranscriptionsResp
}

type CreateAudioTranscriptionsResp struct {
	baseModel
	Data AudioTranscriptionsData `json:"data"`
}

type AudioTranscriptionsData struct {
	Text string `json:"text"`
}

type audioTranscriptions struct {
	core *core
}

func newAudioTranscriptions(core *core) *audioTranscriptions {
	return &audioTranscriptions{core: core}
}
