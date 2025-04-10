package coze

import (
	"context"
	"io"
)

func (r *audioTranscription) Create(ctx context.Context, req *AudioSpeechTranscriptionsReq) (*CreateAudioTranscriptionResp, error) {
	uri := "/v1/audio/transcriptions"
	resp := &CreateAudioTranscriptionResp{}
	if err := r.core.UploadFile(ctx, uri, req.Audio, req.Filename, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type AudioSpeechTranscriptionsReq struct {
	Filename string    `json:"filename"`
	Audio    io.Reader `json:"audio"`
}

type CreateAudioTranscriptionResp struct {
	baseResponse
	Data AudioTranscriptionsData `json:"data"`
}

type AudioTranscriptionsData struct {
	Text string `json:"text"`
}

type audioTranscription struct {
	core *core
}

func newTranscription(core *core) *audioTranscription {
	return &audioTranscription{core: core}
}
