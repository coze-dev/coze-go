package coze

import (
	"context"
	"net/http"
)

// WorkflowMode represents the workflow mode
type WorkflowMode string

const (
	// WorkflowModeWorkflow represents workflow mode
	WorkflowModeWorkflow WorkflowMode = "workflow"
	// WorkflowModeChatflow represents chatflow mode
	WorkflowModeChatflow WorkflowMode = "chatflow"
)

// PublishStatus represents the publish status
type PublishStatus string

const (
	// PublishStatusPublished represents published status
	PublishStatusPublished PublishStatus = "published"
	// PublishStatusUnpublished represents unpublished status
	PublishStatusUnpublished PublishStatus = "unpublished"
)

// WorkflowInfo represents workflow information
type WorkflowInfo struct {
	baseModel
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
	Description  string `json:"description"`
	IconURL      string `json:"icon_url"`
	AppID        string `json:"app_id"`
}

// ListWorkflowsReq represents the request to list workflows
type ListWorkflowsReq struct {
	WorkspaceID   string         `query:"workspace_id" json:"-"`
	WorkflowMode  *WorkflowMode  `query:"workflow_mode" json:"-"`
	AppID         string         `query:"app_id" json:"-"`
	PublishStatus *PublishStatus `query:"publish_status" json:"-"`
	PageNum       int            `query:"page_num" json:"-"`
	PageSize      int            `query:"page_size" json:"-"`
}

// listWorkflowsResp represents the response from listing workflows
type listWorkflowsResp struct {
	baseResponse
	Data *listWorkflowsData `json:"data"`
}

// listWorkflowsData represents the data part of list workflows response
type listWorkflowsData struct {
	Items   []*WorkflowInfo `json:"items"`
	HasMore bool            `json:"has_more"`
}

type workflows struct {
	Runs *workflowRuns
	Chat *workflowsChat
	core *core
}

func newWorkflows(core *core) *workflows {
	return &workflows{
		Runs: newWorkflowRun(core),
		Chat: newWorkflowsChat(core),
		core: core,
	}
}

// List lists workflows with optional filters and pagination
func (w *workflows) List(ctx context.Context, req *ListWorkflowsReq) (NumberPaged[WorkflowInfo], error) {
	if req == nil {
		req = &ListWorkflowsReq{}
	}

	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	fetcher := func(request *pageRequest) (*pageResponse[WorkflowInfo], error) {
		rawReq := &RawRequestReq{
			Method: http.MethodGet,
			URL:    "/v1/workflows",
			Body: &ListWorkflowsReq{
				WorkspaceID:   req.WorkspaceID,
				WorkflowMode:  req.WorkflowMode,
				AppID:         req.AppID,
				PublishStatus: req.PublishStatus,
				PageNum:       request.PageNum,
				PageSize:      request.PageSize,
			},
		}

		resp := new(listWorkflowsResp)
		err := w.core.rawRequest(ctx, rawReq, resp)
		if err != nil {
			return nil, err
		}

		return &pageResponse[WorkflowInfo]{
			response: resp.HTTPResponse,
			HasMore:  resp.Data.HasMore,
			Data:     resp.Data.Items,
			LogID:    resp.HTTPResponse.LogID(),
		}, nil
	}

	return NewNumberPaged[WorkflowInfo](fetcher, req.PageSize, req.PageNum)
}
