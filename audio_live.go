package coze

import (
	"context"
	"net/http"
)

func (r *audioLive) Retrieve(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    buildSwaggerOperationURL("/v1/audio/live/{live_id}", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type audioLive struct {
	core *core
}

func newAudioLive(core *core) *audioLive {
	return &audioLive{core: core}
}
