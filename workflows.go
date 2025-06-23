package coze

import (
	"context"
	"net/http"
	"time"
)

// WorkflowMode represents the mode of a workflow
type WorkflowMode string

const (
	WorkflowModeManual WorkflowMode = "manual"
	WorkflowModeAuto   WorkflowMode = "auto"
)

// WorkflowInfo contains information about a workflow
 type WorkflowInfo struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Mode        WorkflowMode `json:"mode"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
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

// List gets paginated workflows with optional filters
func (w *workflows) List(ctx context.Context, req *ListWorkflowsRequest) (NumberPaged[WorkflowInfo], error) {
	return NewNumberPaged(func(pageReq *pageRequest) (*pageResponse[WorkflowInfo], error) {
		resp := &struct {
			pageResponse[WorkflowInfo]
		}{}
		if err := w.core.rawRequest(ctx, &RawRequestReq{
			Method: http.MethodGet,
			URL:    "/v1/workflows",
			Body: struct {
				WorkspaceID string       `query:"workspace_id"`
				Mode        WorkflowMode `query:"mode"`
				PageNum     int          `query:"page_num"`
				PageSize    int          `query:"page_size"`
			}{req.WorkspaceID, req.Mode, pageReq.PageNum, pageReq.PageSize},
		}, resp); err != nil {
			return nil, err
		}
		return &resp.pageResponse, nil
	}, req.PageSize, req.PageNum)
}

// ListWorkflowsRequest contains parameters for listing workflows
 type ListWorkflowsRequest struct {
	WorkspaceID string       `json:"workspace_id,omitempty"`
	Mode        WorkflowMode `json:"mode,omitempty"`
	PageSize    int          `json:"page_size,omitempty"`
	PageNum     int          `json:"page_num,omitempty"`
}
