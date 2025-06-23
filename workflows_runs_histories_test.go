package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowRunsHistories(t *testing.T) {
	as := assert.New(t)
	t.Run("retrieve workflow run history success", func(t *testing.T) {
		histories := newWorkflowRunsHistories(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workflows/workflow1/run_histories/exec1", req.URL.Path)
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
		})))
		resp, err := histories.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesReq{
			WorkflowID: "workflow1",
			ExecuteID:  "exec1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Len(resp.Histories, 1)

		history := resp.Histories[0]
		as.Equal("exec1", history.ExecuteID)
		as.Equal(WorkflowExecuteStatusSuccess, history.ExecuteStatus)
		as.Equal("bot1", history.BotID)
		as.Equal("1024", history.ConnectorID)
		as.Equal("user1", history.ConnectorUid)
		as.Equal(WorkflowRunModeStreaming, history.RunMode)
		as.Equal("log1", history.LogID)
		as.Equal(1234567890, history.CreateTime)
		as.Equal(1234567891, history.UpdateTime)
		as.Equal(`{"result": "success"}`, history.Output)
		as.Equal("0", history.ErrorCode)
		as.Empty(history.ErrorMessage)
		as.Equal("https://debug.example.com", history.DebugURL)
	})

	t.Run("retrieve workflow run history with error", func(t *testing.T) {
		histories := newWorkflowRunsHistories(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := histories.Retrieve(context.Background(), &RetrieveWorkflowsRunsHistoriesReq{
			WorkflowID: "invalid_workflow",
			ExecuteID:  "invalid_exec",
		})
		as.NotNil(err)
	})
}

func TestWorkflowRunMode(t *testing.T) {
	as := assert.New(t)
	t.Run("workflow run mode constants", func(t *testing.T) {
		as.Equal(WorkflowRunMode(0), WorkflowRunModeSynchronous)
		as.Equal(WorkflowRunMode(1), WorkflowRunModeStreaming)
		as.Equal(WorkflowRunMode(2), WorkflowRunModeAsynchronous)
	})
	t.Run("workflow execute status constants", func(t *testing.T) {
		as.Equal(WorkflowExecuteStatus("Success"), WorkflowExecuteStatusSuccess)
		as.Equal(WorkflowExecuteStatus("Running"), WorkflowExecuteStatusRunning)
		as.Equal(WorkflowExecuteStatus("Fail"), WorkflowExecuteStatusFail)
	})
}
