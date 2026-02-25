package coze

import (
	"context"
	"net/http"
)

// Upload 上传文件
func (r *files) Upload(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    buildSwaggerOperationURL("/v1/files/upload", req.PathParams, req.QueryParams),
		Body:   req.Body,
		IsFile: true,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

// Retrieve 查看文件详情
func (r *files) Retrieve(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    buildSwaggerOperationURL("/v1/files/retrieve", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type files struct {
	core *core
}

func newFiles(core *core) *files {
	return &files{core: core}
}
