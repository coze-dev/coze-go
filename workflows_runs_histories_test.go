package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowRunsHistories(t *testing.T) {
	// Test Retrieve method
	t.Run("Retrieve workflow run history success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/workflows/workflow1/run_histories/exec1", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &retrieveWorkflowRunsHistoriesResp{
					RetrieveWorkflowRunsHistoriesResp: &RetrieveWorkflowRunsHistoriesResp{
						Histories: []*WorkflowRunHistory{
							{
								ExecuteID:     "exec1",
								ExecuteStatus: WorkflowExecuteStatusSuccess,
								BotID:         "bot1",
								ConnectorID:   "1024",
								ConnectorUid:  "user1",
								RunMode:       WorkflowRunModeStreaming,
								LogID:         "log1",
								CreateTime:    1234567890,
								UpdateTime:    1234567891,
								Output:        `{"result": "success"}`,
								ErrorCode:     "0",
								ErrorMessage:  "",
								DebugURL:      "https://debug.example.com",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		histories := newWorkflowRunsHistories(core)

		resp, err := histories.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesReq{
			WorkflowID: "workflow1",
			ExecuteID:  "exec1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		require.Len(t, resp.Histories, 1)

		history := resp.Histories[0]
		assert.Equal(t, "exec1", history.ExecuteID)
		assert.Equal(t, WorkflowExecuteStatusSuccess, history.ExecuteStatus)
		assert.Equal(t, "bot1", history.BotID)
		assert.Equal(t, "1024", history.ConnectorID)
		assert.Equal(t, "user1", history.ConnectorUid)
		assert.Equal(t, WorkflowRunModeStreaming, history.RunMode)
		assert.Equal(t, "log1", history.LogID)
		assert.Equal(t, 1234567890, history.CreateTime)
		assert.Equal(t, 1234567891, history.UpdateTime)
		assert.Equal(t, `{"result": "success"}`, history.Output)
		assert.Equal(t, "0", history.ErrorCode)
		assert.Empty(t, history.ErrorMessage)
		assert.Equal(t, "https://debug.example.com", history.DebugURL)
	})

	// Test Retrieve method with error
	t.Run("Retrieve workflow run history with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		histories := newWorkflowRunsHistories(core)

		resp, err := histories.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesReq{
			WorkflowID: "invalid_workflow",
			ExecuteID:  "invalid_exec",
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestWorkflowRunMode(t *testing.T) {
	t.Run("WorkflowRunMode constants", func(t *testing.T) {
		assert.Equal(t, WorkflowRunMode(0), WorkflowRunModeSynchronous)
		assert.Equal(t, WorkflowRunMode(1), WorkflowRunModeStreaming)
		assert.Equal(t, WorkflowRunMode(2), WorkflowRunModeAsynchronous)
	})
}

func TestWorkflowExecuteStatus(t *testing.T) {
	t.Run("WorkflowExecuteStatus constants", func(t *testing.T) {
		assert.Equal(t, WorkflowExecuteStatus("Success"), WorkflowExecuteStatusSuccess)
		assert.Equal(t, WorkflowExecuteStatus("Running"), WorkflowExecuteStatusRunning)
		assert.Equal(t, WorkflowExecuteStatus("Fail"), WorkflowExecuteStatusFail)
	})
}
