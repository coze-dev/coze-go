package coze

import (
	"context"
	"net/http"
)

func (r *workflowRunsHistories) Retrieve(ctx context.Context, req *RetrieveWorkflowsRunsHistoriesReq) (*RetrieveWorkflowRunsHistoriesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/workflows/:workflow_id/run_histories/:execute_id",
		Body:   req,
	}
	response := new(retrieveWorkflowRunsHistoriesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.RetrieveWorkflowRunsHistoriesResp, err
}

// WorkflowRunMode represents how the workflow runs
type WorkflowRunMode int

const (
	// WorkflowRunModeSynchronous Synchronous operation.
	WorkflowRunModeSynchronous WorkflowRunMode = 0

	// WorkflowRunModeStreaming Streaming operation.
	WorkflowRunModeStreaming WorkflowRunMode = 1

	// WorkflowRunModeAsynchronous Asynchronous operation.
	WorkflowRunModeAsynchronous WorkflowRunMode = 2
)

// WorkflowExecuteStatus represents the execution status of a workflow
type WorkflowExecuteStatus string

const (
	// WorkflowExecuteStatusSuccess Execution succeeded.
	WorkflowExecuteStatusSuccess WorkflowExecuteStatus = "Success"

	// WorkflowExecuteStatusRunning Execution in progress.
	WorkflowExecuteStatusRunning WorkflowExecuteStatus = "Running"

	// WorkflowExecuteStatusFail Execution failed.
	WorkflowExecuteStatusFail WorkflowExecuteStatus = "Fail"
)

// RetrieveWorkflowsRunsHistoriesReq represents request for retrieving workflow runs history
type RetrieveWorkflowsRunsHistoriesReq struct {
	// The ID of the workflow.
	ExecuteID string `path:"execute_id" json:"-"`

	// The ID of the workflow async execute.
	WorkflowID string `path:"workflow_id" json:"-"`
}

// RunWorkflowsResp represents response for running workflow
type RunWorkflowsResp struct {
	baseModel
	// Execution ID of asynchronous execution.
	ExecuteID string `json:"execute_id,omitempty"`

	// Workflow execution result.
	Data string `json:"data,omitempty"`

	DebugURL string `json:"debug_url,omitempty"`
	Token    int    `json:"token,omitempty"`
	Cost     string `json:"cost,omitempty"`
}

// RetrieveWorkflowRunsHistoriesResp represents response for retrieving workflow runs history
type RetrieveWorkflowRunsHistoriesResp struct {
	baseModel
	Histories []*WorkflowRunHistory `json:"data"`
}

// WorkflowRunHistory represents the history of a workflow runs
type WorkflowRunHistory struct {
	// The ID of execute.
	ExecuteID string `json:"execute_id"`

	// Execute status: success: Execution succeeded. running: Execution in progress. fail: Execution failed.
	ExecuteStatus WorkflowExecuteStatus `json:"execute_status"`

	// The Bot ID specified when executing the workflow. Returns 0 if no Bot ID was specified.
	BotID string `json:"bot_id"`

	// The release connector ID of the agent. By default, only the Agent as API connector is
	// displayed, and the connector ID is 1024.
	ConnectorID string `json:"connector_id"`

	// User ID, the user_id specified by the ext field when executing the workflow. If not specified,
	// the token applicant's button ID is returned.
	ConnectorUid string `json:"connector_uid"`

	// How the workflow runs: 0: Synchronous operation. 1: Streaming operation. 2: Asynchronous operation.
	RunMode WorkflowRunMode `json:"run_mode"`

	// The Log ID of the asynchronously running workflow. If the workflow is executed abnormally, you
	// can contact the service team to troubleshoot the problem through the Log ID.
	LogID string `json:"logid"`

	// The start time of the workflow, in Unix time timestamp format, in seconds.
	CreateTime int `json:"create_time"`

	// The workflow resume running time, in Unix time timestamp format, in seconds.
	UpdateTime int `json:"update_time"`

	// The output of the workflow is usually a JSON serialized string, but it may also be a non-JSON
	// structured string.
	Output string `json:"output"`

	// The status of each output node execution.
	NodeExecuteStatus map[string]*WorkflowRunHistoryNodeExecuteStatus `json:"node_execute_status"`

	// Status code. 0 represents a successful API call. Other values indicate that the call has
	// failed. You can determine the detailed reason for the error through the error_message field.
	ErrorCode string `json:"error_code"`

	// Status message. You can get detailed error information when the API call fails.
	ErrorMessage string `json:"error_message"`

	// Workflow trial runs debugging page. Visit this page to view the running results, input and
	// output information of each workflow node.
	DebugURL string `json:"debug_url"`
}

// NodeExecuteStatus represents the status of a node execution
type WorkflowRunHistoryNodeExecuteStatus struct {
	// The ID of the node.
	NodeID string `json:"node_id"`
	// Whether the node is finished.
	IsFinish bool `json:"is_finish"`
	// The loop index of the node.
	LoopIndex *int `json:"loop_index"`
	// The batch index of the node.
	BatchIndex *int `json:"batch_index"`
	// The update time of the node.
	UpdateTime int `json:"update_time"`
	// The ID of the sub-execute.
	SubExecuteID *string `json:"sub_execute_id"`
	// The UUID of the node execution.
	NodeExecuteUUID string `json:"node_execute_uuid"`
}

type runWorkflowsResp struct {
	baseResponse
	*RunWorkflowsResp
	HTTPResponse *http.Response `json:"-"`
}

type retrieveWorkflowRunsHistoriesResp struct {
	baseResponse
	*RetrieveWorkflowRunsHistoriesResp
}

type workflowRunsHistories struct {
	core *core

	ExecuteNodes *workflowsRunsHistoriesExecuteNodes
}

func newWorkflowRunsHistories(core *core) *workflowRunsHistories {
	return &workflowRunsHistories{
		core:         core,
		ExecuteNodes: newWorkflowsRunsHistoriesExecuteNodes(core),
	}
}
