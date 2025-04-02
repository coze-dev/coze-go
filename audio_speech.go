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

func (r *audioSpeech) Transcription(ctx context.Context, reader io.Reader, filename string) (*CreateAudioTranscriptionResp, error) {
	uri := "/v1/audio/transcriptions"
	resp := &CreateAudioTranscriptionResp{}
	if err := r.core.UploadFile(ctx, uri, reader, filename, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
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
}

// CreateAudioSpeechResp represents the response for creating speech
type CreateAudioSpeechResp struct {
	baseResponse
	// TODO 没有 json tag？
	Data io.ReadCloser
}

type CreateAudioTranscriptionResp struct {
	baseResponse
	Data AudioTranscriptionsData `json:"data"`
}

type AudioTranscriptionsData struct {
	Text string `json:"text"`
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
