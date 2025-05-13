package coze

import (
	"context"
	"io"
	"net/http"
	"os"
)

func (r *audioSpeech) Create(ctx context.Context, req *CreateAudioSpeechReq) (*CreateAudioSpeechResp, error) {
	uri := "/v1/audio/speech"
	resp, err := r.core.RawRequest(ctx, http.MethodPost, uri, req)
	if err != nil {
		return nil, err
	}
	res := &CreateAudioSpeechResp{
		Data: resp.Body,
	}
	res.SetHTTPResponse(newHTTPResponse(resp))
	return res, nil
}

type audioSpeech struct {
	core *core
}

func newSpeech(core *core) *audioSpeech {
	return &audioSpeech{core: core}
}

// CreateAudioSpeechReq represents the request for creating speech
type CreateAudioSpeechReq struct {
	Input          string       `json:"input"`
	VoiceID        string       `json:"voice_id"`
	ResponseFormat *AudioFormat `json:"response_format"`
	Speed          *float32     `json:"speed"`
	SampleRate     *int         `json:"sample_rate"`
}

// CreateAudioSpeechResp represents the response for creating speech
type CreateAudioSpeechResp struct {
	baseResponse
	// TODO 没有 json tag？
	Data io.ReadCloser
}

func (c *CreateAudioSpeechResp) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	defer c.Data.Close()

	_, err = io.Copy(file, c.Data)
	return err
}
