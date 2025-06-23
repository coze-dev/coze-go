package coze

import (
	"context"
	"net/http"
)

// Retrieve retrieves the output of a node execution
func (r *workflowsRunsHistoriesExecuteNodes) Retrieve(ctx context.Context, req *RetrieveWorkflowsRunsHistoriesExecuteNodesReq) (*RetrieveWorkflowRunsHistoriesExecuteNodesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/workflows/:workflow_id/run_histories/:execute_id/execute_nodes/:node_execute_uuid",
		Body:   req,
	}
	response := new(retrieveWorkflowRunsHistoriesExecuteNodeResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// RetrieveWorkflowsRunsHistoriesExecuteNodesReq query output node execution result
type RetrieveWorkflowsRunsHistoriesExecuteNodesReq struct {
	// The ID of the workflow async execute.
	WorkflowID string `path:"workflow_id" json:"-"`

	// The ID of the workflow.
	ExecuteID string `path:"execute_id" json:"-"`

	// The ID of the node execute.
	NodeExecuteUUID string `path:"node_execute_uuid" json:"-"`
}

// RetrieveWorkflowRunsHistoriesExecuteNodesResp allows you to retrieve the output of a node execution
type RetrieveWorkflowRunsHistoriesExecuteNodesResp struct {
	baseModel
	// The node is finished.
	IsFinish bool `json:"is_finish"`
	// The node output.
	NodeOutput string `json:"node_output"`
}

type retrieveWorkflowRunsHistoriesExecuteNodeResp struct {
	baseResponse
	Data *RetrieveWorkflowRunsHistoriesExecuteNodesResp `json:"data"`
}

type workflowsRunsHistoriesExecuteNodes struct {
	core *core
}

func newWorkflowsRunsHistoriesExecuteNodes(core *core) *workflowsRunsHistoriesExecuteNodes {
	return &workflowsRunsHistoriesExecuteNodes{core: core}
}
