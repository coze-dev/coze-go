package coze

import (
	"context"
	"net/http"
)

// Stream 执行对话流
func (r *workflowsChat) Stream(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/workflows/chat", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

// Create 执行对话流
func (r *workflowsChat) Create(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/workflows/chat", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type workflowsChat struct {
	core *core
}

func newWorkflowsChat(core *core) *workflowsChat {
	return &workflowsChat{core: core}
}
