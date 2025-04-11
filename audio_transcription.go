package coze

import (
	"context"
	"io"
)

func (r *audioTranscriptions) Create(ctx context.Context, req *AudioSpeechTranscriptionsReq) (*CreateAudioTranscriptionsResp, error) {
	uri := "/v1/audio/transcriptions"
	resp := &CreateAudioTranscriptionsResp{}
	if err := r.core.UploadFile(ctx, uri, req.Audio, req.Filename, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type AudioSpeechTranscriptionsReq struct {
	Filename string    `json:"filename"`
	Audio    io.Reader `json:"audio"`
}

type CreateAudioTranscriptionsResp struct {
	baseResponse
	Data AudioTranscriptionsData `json:"data"`
}

type AudioTranscriptionsData struct {
	Text string `json:"text"`
}

type audioTranscriptions struct {
	core *core
}

func newTranscriptions(core *core) *audioTranscriptions {
	return &audioTranscriptions{core: core}
}
