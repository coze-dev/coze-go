package coze

import (
	"context"
	"io"
	"net/http"
	"os"
)

func (r *audioSpeech) Create(ctx context.Context, req *CreateAudioSpeechReq) (*CreateAudioSpeechResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/speech",
		Body:   req,
	}
	response := new(createAudioSpeechResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
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
	baseModel
	Data io.ReadCloser
}

type createAudioSpeechResp struct {
	baseResponse
	Data *CreateAudioSpeechResp
}

func (r *createAudioSpeechResp) SetReader(file io.ReadCloser) {
	if r.Data == nil {
		r.Data = &CreateAudioSpeechResp{}
	}
	r.Data.Data = file
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

type audioSpeech struct {
	core *core
}

func newSpeech(core *core) *audioSpeech {
	return &audioSpeech{core: core}
}
