package coze

import (
	"context"
	"net/http"
)

// List 查看对话消息详情
func (r *chatMessages) List(ctx context.Context, req *SwaggerOperationRequest) (*SwaggerOperationResponse, error) {
	if req == nil {
		req = &SwaggerOperationRequest{}
	}
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    buildSwaggerOperationURL("/v3/chat/message/list", req.PathParams, req.QueryParams),
		Body:   req.Body,
	}
	response := new(SwaggerOperationResponse)
	err := r.core.rawRequest(ctx, request, response)
	return response, err
}

type chatMessages struct {
	core *core
}

func newChatMessages(core *core) *chatMessages {
	return &chatMessages{core: core}
}
