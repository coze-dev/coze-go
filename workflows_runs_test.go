package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowRuns(t *testing.T) {
	// Test Create method
	t.Run("Create workflow run success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/workflow/run", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &runWorkflowsResp{
					RunWorkflowsResp: &RunWorkflowsResp{
						ExecuteID: "exec1",
						Data:      `{"result": "success"}`,
						DebugURL:  "https://debug.example.com",
						Token:     100,
						Cost:      "0.1",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workflowRuns := newWorkflowRun(core)

		resp, err := workflowRuns.Create(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
			Parameters: map[string]any{
				"param1": "value1",
			},
			BotID:   "bot1",
			IsAsync: true,
			AppID:   "app1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "exec1", resp.ExecuteID)
		assert.Equal(t, `{"result": "success"}`, resp.Data)
		assert.Equal(t, "https://debug.example.com", resp.DebugURL)
		assert.Equal(t, 100, resp.Token)
		assert.Equal(t, "0.1", resp.Cost)
	})

	// Test Stream method
	t.Run("Stream workflow run success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/workflow/stream_run", req.URL.Path)

				// Return mock streaming response
				return mockStreamResponse(`id:0
event:Message
data:{"content":"Hello","node_title":"Start","node_seq_id":"0","node_is_finish":false}

id:1
event:Message
data:{"content":"World","node_title":"End","node_seq_id":"1","node_is_finish":true}

id:2
event:Done
data:{"debug_url":"https://www.coze.cn/work_flow?***"}
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workflowRuns := newWorkflowRun(core)

		reader, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
			Parameters: map[string]any{
				"param1": "value1",
			},
		})

		require.NoError(t, err)
		defer reader.Close()

		// Read first message event
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, 0, event.ID)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "Hello", event.Message.Content)
		assert.Equal(t, "Start", event.Message.NodeTitle)
		assert.Equal(t, "0", event.Message.NodeSeqID)
		assert.False(t, event.Message.NodeIsFinish)

		// Read second message event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, 1, event.ID)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "World", event.Message.Content)
		assert.Equal(t, "End", event.Message.NodeTitle)
		assert.Equal(t, "1", event.Message.NodeSeqID)
		assert.True(t, event.Message.NodeIsFinish)

		// Read done event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, 2, event.ID)
		assert.Equal(t, WorkflowEventTypeDone, event.Event)
		assert.Equal(t, "https://www.coze.cn/work_flow?***", event.DebugURL.URL)
		assert.True(t, event.IsDone())
	})

	// Test Resume method
	t.Run("Resume workflow run success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/workflow/stream_resume", req.URL.Path)

				// Return mock streaming response
				return mockStreamResponse(`id:0
event:Message
data:{"content":"Resumed","node_title":"Resume","node_seq_id":"0","node_is_finish":true}

id:1
event:Done
data:{"debug_url":"https://www.coze.cn/work_flow?***"}
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workflowRuns := newWorkflowRun(core)

		reader, err := workflowRuns.Resume(context.Background(), &ResumeRunWorkflowsReq{
			WorkflowID:    "workflow1",
			EventID:       "event1",
			ResumeData:    "data1",
			InterruptType: 1,
		})

		require.NoError(t, err)
		defer reader.Close()

		// Read message event
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, 0, event.ID)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "Resumed", event.Message.Content)
		assert.Equal(t, "Resume", event.Message.NodeTitle)
		assert.Equal(t, "0", event.Message.NodeSeqID)
		assert.True(t, event.Message.NodeIsFinish)

		// Read done event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, 1, event.ID)
		assert.Equal(t, WorkflowEventTypeDone, event.Event)
		assert.Equal(t, "https://www.coze.cn/work_flow?***", event.DebugURL.URL)
		assert.True(t, event.IsDone())
	})

	// Test error event parsing
	t.Run("Parse error event", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockStreamResponse(`id:0
event:Error
data:{"error_code":400,"error_message":"Bad Request"}
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workflowRuns := newWorkflowRun(core)

		reader, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
		})

		require.NoError(t, err)
		defer reader.Close()

		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeError, event.Event)
		assert.Equal(t, 400, event.Error.ErrorCode)
		assert.Equal(t, "Bad Request", event.Error.ErrorMessage)
	})

	// Test interrupt event parsing
	t.Run("Parse interrupt event", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockStreamResponse(`id:0
event:Interrupt
data:{"interrupt_data":{"event_id":"event1","type":1},"node_title":"Question"}
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workflowRuns := newWorkflowRun(core)

		reader, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
		})

		require.NoError(t, err)
		defer reader.Close()

		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeInterrupt, event.Event)
		assert.Equal(t, "event1", event.Interrupt.InterruptData.EventID)
		assert.Equal(t, 1, event.Interrupt.InterruptData.Type)
		assert.Equal(t, "Question", event.Interrupt.NodeTitle)
	})
}

func TestWorkflowEventParsing(t *testing.T) {
	t.Run("ParseWorkflowEventError", func(t *testing.T) {
		data := `{"error_code":400,"error_message":"Bad Request"}`
		err, parseErr := ParseWorkflowEventError(data)
		require.NoError(t, parseErr)
		assert.Equal(t, 400, err.ErrorCode)
		assert.Equal(t, "Bad Request", err.ErrorMessage)
	})

	t.Run("ParseWorkflowEventInterrupt", func(t *testing.T) {
		data := `{"interrupt_data":{"event_id":"event1","type":1},"node_title":"Question"}`
		interrupt, parseErr := ParseWorkflowEventInterrupt(data)
		require.NoError(t, parseErr)
		assert.Equal(t, "event1", interrupt.InterruptData.EventID)
		assert.Equal(t, 1, interrupt.InterruptData.Type)
		assert.Equal(t, "Question", interrupt.NodeTitle)
	})

	t.Run("Invalid JSON parsing", func(t *testing.T) {
		_, err := ParseWorkflowEventError("invalid json")
		assert.Error(t, err)

		_, err = ParseWorkflowEventInterrupt("invalid json")
		assert.Error(t, err)
	})
}
