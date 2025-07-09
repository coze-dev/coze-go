package coze

import (
	"context"
	"net/http"
)

func (r *workflows) List(ctx context.Context, req *ListWorkflowReq) (NumberPaged[WorkflowInfo], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[WorkflowInfo], error) {
			resp := new(listWorkflowResp)
			err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/workflows",
				Body:   req.toReq(request),
			}, resp)
			if err != nil {
				return nil, err
			}
			return &pageResponse[WorkflowInfo]{
				response: resp.HTTPResponse,
				HasMore:  resp.Data.HasMore,
				Data:     resp.Data.Items,
				LogID:    resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type ListWorkflowReq struct {
	WorkspaceID   *string        `query:"workspace_id" json:"-"`
	WorkflowMode  *WorkflowMode  `query:"workflow_mode" json:"-"`
	AppID         *string        `query:"app_id" json:"-"`
	PublishStatus *PublishStatus `query:"publish_status" json:"-"`
	PageNum       int            `query:"page_num" json:"-"`
	PageSize      int            `query:"page_size" json:"-"`
}

type WorkflowMode string

const (
	WorkflowModeWorkflow WorkflowMode = "workflow"
	WorkflowModeChatflow WorkflowMode = "chatflow"
)

type WorkflowInfo struct {
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
	Description  string `json:"description"`
	IconURL      string `json:"icon_url"`
	AppID        string `json:"app_id"`
}

type PublishStatus string

const (
	// PublishStatusALL 所有智能体，且数据为最新草稿版本
	PublishStatusALL PublishStatus = "all"
	// PublishStatusPublishedOnline 已发布智能体的最新线上版本
	PublishStatusPublishedOnline PublishStatus = "published_online"
	// PublishStatusPublishedDraft 已发布的最新草稿版本
	PublishStatusPublishedDraft PublishStatus = "published_draft"
	// PublishStatusUnpublishedDraft 未发布的最新草稿版本
	PublishStatusUnpublishedDraft PublishStatus = "unpublished_draft"
)

type ListWorkflowResp struct {
	HasMore bool            `json:"has_more"`
	Items   []*WorkflowInfo `json:"items"`
}

type listWorkflowResp struct {
	baseResponse
	Data *ListWorkflowResp `json:"data"`
}

func (r ListWorkflowReq) toReq(request *pageRequest) *ListWorkflowReq {
	return &ListWorkflowReq{
		WorkspaceID:   r.WorkspaceID,
		WorkflowMode:  r.WorkflowMode,
		AppID:         r.AppID,
		PublishStatus: r.PublishStatus,
		PageNum:       request.PageNum,
		PageSize:      request.PageSize,
	}
}

type workflows struct {
	core *core
	Runs *workflowRuns
	Chat *workflowsChat
}

func newWorkflows(core *core) *workflows {
	return &workflows{
		core: core,
		Runs: newWorkflowRun(core),
		Chat: newWorkflowsChat(core),
	}
}
