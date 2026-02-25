package coze

import (
	"context"
	"net/http"
)

// Create 语音识别
func (r *audioTranscriptions) Create(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/audio/transcriptions", req.PathParams, req.QueryParams),
		Body:   req.Body,
		IsFile: true,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type audioTranscriptions struct {
	core *core
}

func newAudioTranscriptions(core *core) *audioTranscriptions {
	return &audioTranscriptions{core: core}
}
