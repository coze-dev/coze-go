package coze

import (
	"context"
	"fmt"
	"net/http"
)

// RetrieveNode retrieves the output of a node execution
func (r *workflowRunsHistories) RetrieveNode(ctx context.Context, req *RetrieveWorkflowsRunsHistoriesNodeReq) (*RetrieveWorkflowRunsHistoriesNodeResp, error) {
	method := http.MethodGet
	uri := fmt.Sprintf("/v1/workflows/%s/run_histories/%s/execute_nodes/%s", req.WorkflowID, req.ExecuteID, req.NodeExecuteUUID)
	resp := &retrieveWorkflowRunsHistoriesNodeResp{}
	err := r.core.Request(ctx, method, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	resp.RetrieveWorkflowRunsHistoriesNodeResp.setHTTPResponse(resp.HTTPResponse)
	return resp.RetrieveWorkflowRunsHistoriesNodeResp, nil
}

// RetrieveWorkflowsRunsHistoriesNodeReq query output node execution result
type RetrieveWorkflowsRunsHistoriesNodeReq struct {
	// The ID of the workflow.
	ExecuteID string `json:"execute_id"`

	// The ID of the workflow async execute.
	WorkflowID string `json:"workflow_id"`

	// The ID of the node execute.
	NodeExecuteUUID string `json:"node_execute_uuid"`
}

type retrieveWorkflowRunsHistoriesNodeResp struct {
	baseResponse
	*RetrieveWorkflowRunsHistoriesNodeResp
}

// NodeResult The node result.
type NodeResult struct {
	// The node is finished.
	IsFinish bool `json:"is_finish"`
	// The node output.
	NodeOutput string `json:"node_output"`
}

// RetrieveWorkflowRunsHistoriesNodeResp allows you to retrieve the output of a node execution
type RetrieveWorkflowRunsHistoriesNodeResp struct {
	baseModel
	NodeResult `json:"data"`
}
