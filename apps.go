package coze

import (
	"context"
	"net/http"
)

// List 查看应用列表
func (r *apps) List(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    buildSwaggerOperationURL("/v1/apps", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type apps struct {
	core *core
}

func newApps(core *core) *apps {
	return &apps{core: core}
}
