package coze

import (
	"context"
	"fmt"
	"net/http"
)

type workflowsRunsHistoriesExecuteNodes struct {
	core *core
}

func newWorkflowsRunsHistoriesExecuteNodes(core *core) *workflowsRunsHistoriesExecuteNodes {
	return &workflowsRunsHistoriesExecuteNodes{core: core}
}

// Retrieve retrieves the output of a node execution
func (r *workflowsRunsHistoriesExecuteNodes) Retrieve(ctx context.Context, req *RetrieveWorkflowsRunsHistoriesExecuteNodesReq) (*RetrieveWorkflowRunsHistoriesExecuteNodesResp, error) {
	method := http.MethodGet
	uri := fmt.Sprintf("/v1/workflows/%s/run_histories/%s/execute_nodes/%s", req.WorkflowID, req.ExecuteID, req.NodeExecuteUUID)
	resp := &retrieveWorkflowRunsHistoriesExecuteNodeResp{}
	err := r.core.Request(ctx, method, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

// RetrieveWorkflowsRunsHistoriesExecuteNodesReq query output node execution result
type RetrieveWorkflowsRunsHistoriesExecuteNodesReq struct {
	// The ID of the workflow async execute.
	WorkflowID string `json:"workflow_id"`

	// The ID of the workflow.
	ExecuteID string `json:"execute_id"`

	// The ID of the node execute.
	NodeExecuteUUID string `json:"node_execute_uuid"`
}

type retrieveWorkflowRunsHistoriesExecuteNodeResp struct {
	baseResponse
	Data *RetrieveWorkflowRunsHistoriesExecuteNodesResp `json:"data"`
}

// RetrieveWorkflowRunsHistoriesExecuteNodesResp allows you to retrieve the output of a node execution
type RetrieveWorkflowRunsHistoriesExecuteNodesResp struct {
	baseModel
	// The node is finished.
	IsFinish bool `json:"is_finish"`
	// The node output.
	NodeOutput string `json:"node_output"`
}
