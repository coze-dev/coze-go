package coze

import (
	"context"
	"net/http"
)

func (r *users) Me(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    buildSwaggerOperationURL("/v1/users/me", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type users struct {
	core *core
}

func newUsers(core *core) *users {
	return &users{core: core}
}
