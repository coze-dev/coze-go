package coze

import (
	"context"
	"net/http"
)

// Create 语音合成
func (r *audioSpeech) Create(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/audio/speech", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type audioSpeech struct {
	core *core
}

func newAudioSpeech(core *core) *audioSpeech {
	return &audioSpeech{core: core}
}
