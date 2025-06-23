package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowRunsHistoriesExecuteNodes(t *testing.T) {
	as := assert.New(t)
	t.Run("retrieve workflow run history execute node success", func(t *testing.T) {
		executeNodes := newWorkflowsRunsHistoriesExecuteNodes(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workflows/workflow1/run_histories/exec1/execute_nodes/node1", req.URL.Path)
			return mockResponse(http.StatusOK, &retrieveWorkflowRunsHistoriesExecuteNodeResp{
				Data: &RetrieveWorkflowRunsHistoriesExecuteNodesResp{
					IsFinish:   true,
					NodeOutput: `{"result": "success"}`,
				},
			})
		})))
		resp, err := executeNodes.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesExecuteNodesReq{
			WorkflowID:      "workflow1",
			ExecuteID:       "exec1",
			NodeExecuteUUID: "node1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.True(resp.IsFinish)
		as.Equal(`{"result": "success"}`, resp.NodeOutput)
	})

	t.Run("retrieve workflow run history execute node with error", func(t *testing.T) {
		executeNodes := newWorkflowsRunsHistoriesExecuteNodes(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := executeNodes.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesExecuteNodesReq{
			WorkflowID:      "invalid_workflow",
			ExecuteID:       "invalid_exec",
			NodeExecuteUUID: "invalid_node",
		})
		as.NotNil(err)
	})
}
