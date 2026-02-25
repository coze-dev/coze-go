package coze

import (
	"context"
	"net/http"
)

func (r *templates) Duplicate(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/templates/{template_id}/duplicate", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type templates struct {
	core *core
}

func newTemplates(core *core) *templates {
	return &templates{core: core}
}
