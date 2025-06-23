package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowRuns(t *testing.T) {
	as := assert.New(t)
	t.Run("create workflow run success", func(t *testing.T) {
		workflowRuns := newWorkflowRun(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/workflow/run", req.URL.Path)
			return mockResponse(http.StatusOK, &runWorkflowsResp{
				RunWorkflowsResp: &RunWorkflowsResp{
					ExecuteID: "exec1",
					Data:      `{"result": "success"}`,
					DebugURL:  "https://debug.example.com",
					Token:     100,
					Cost:      "0.1",
				},
			})
		})))
		resp, err := workflowRuns.Create(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
			Parameters: map[string]any{
				"param1": "value1",
			},
			BotID:   "bot1",
			IsAsync: true,
			AppID:   "app1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal("exec1", resp.ExecuteID)
		as.Equal(`{"result": "success"}`, resp.Data)
		as.Equal("https://debug.example.com", resp.DebugURL)
		as.Equal(100, resp.Token)
		as.Equal("0.1", resp.Cost)
	})

	t.Run("stream workflow run success", func(t *testing.T) {
		workflowRuns := newWorkflowRun(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/workflow/stream_run", req.URL.Path)
			return mockStreamResponse(`id:0
event:Message
data:{"content":"Hello","node_title":"Start","node_seq_id":"0","node_is_finish":false}

id:1
event:Message
data:{"content":"World","node_title":"End","node_seq_id":"1","node_is_finish":true}

id:2
event:PING
data:{}

id:3
event:invalid
data:{}

id:4
event:Done
data:{"debug_url":"https://www.coze.cn/work_flow?***"}
`)
		})))
		stream, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
			Parameters: map[string]any{
				"param1": "value1",
			},
		})
		as.Nil(err)
		as.NotNil(stream)
		as.NotEmpty(stream.Response().LogID())
		defer stream.Close()

		event, err := stream.Recv()
		as.Nil(err)
		as.Equal(0, event.ID)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("Hello", event.Message.Content)
		as.Equal("Start", event.Message.NodeTitle)
		as.Equal("0", event.Message.NodeSeqID)
		as.False(event.Message.NodeIsFinish)

		// read second message event
		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(1, event.ID)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("World", event.Message.Content)
		as.Equal("End", event.Message.NodeTitle)
		as.Equal("1", event.Message.NodeSeqID)
		as.True(event.Message.NodeIsFinish)

		// read ping event
		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(2, event.ID)
		as.Equal(WorkflowEventTypePing, event.Event)

		// read invalid event
		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(3, event.ID)
		as.Equal(WorkflowEventTypeUnknown, event.Event)

		// read done event
		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(4, event.ID)
		as.Equal(WorkflowEventTypeDone, event.Event)
		as.Equal("https://www.coze.cn/work_flow?***", event.DebugURL.URL)
		as.True(event.IsDone())
	})

	t.Run("resume workflow run success", func(t *testing.T) {
		workflowRuns := newWorkflowRun(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/workflow/stream_resume", req.URL.Path)
			return mockStreamResponse(`id:0
event:Message
data:{"content":"Resumed","node_title":"Resume","node_seq_id":"0","node_is_finish":true}

id:1
event:Done
data:{"debug_url":"https://www.coze.cn/work_flow?***"}
`)
		})))
		stream, err := workflowRuns.Resume(context.Background(), &ResumeRunWorkflowsReq{
			WorkflowID:    "workflow1",
			EventID:       "event1",
			ResumeData:    "data1",
			InterruptType: 1,
		})
		as.Nil(err)
		as.NotNil(stream)
		as.NotEmpty(stream.Response().LogID())
		defer stream.Close()

		// read message event
		event, err := stream.Recv()
		as.Nil(err)
		as.Equal(0, event.ID)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("Resumed", event.Message.Content)
		as.Equal("Resume", event.Message.NodeTitle)
		as.Equal("0", event.Message.NodeSeqID)
		as.True(event.Message.NodeIsFinish)

		// Read done event
		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(1, event.ID)
		as.Equal(WorkflowEventTypeDone, event.Event)
		as.Equal("https://www.coze.cn/work_flow?***", event.DebugURL.URL)
		as.True(event.IsDone())
	})

	t.Run("parse error event", func(t *testing.T) {
		workflowRuns := newWorkflowRun(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockStreamResponse(`id:0
event:Error
data:{"error_code":400,"error_message":"Bad Request"}
`)
		})))
		stream, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
		})
		as.Nil(err)
		as.NotNil(stream)
		as.NotEmpty(stream.Response().LogID())
		defer stream.Close()

		event, err := stream.Recv()
		as.Nil(err)
		as.Equal(WorkflowEventTypeError, event.Event)
		as.Equal(400, event.Error.ErrorCode)
		as.Equal("Bad Request", event.Error.ErrorMessage)
	})

	t.Run("parse interrupt event", func(t *testing.T) {
		workflowRuns := newWorkflowRun(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockStreamResponse(`id:0
event:Interrupt
data:{"interrupt_data":{"event_id":"event1","type":1},"node_title":"Question"}
`)
		})))
		stream, err := workflowRuns.Stream(context.Background(), &RunWorkflowsReq{
			WorkflowID: "workflow1",
		})
		as.Nil(err)
		as.NotNil(stream)
		as.NotEmpty(stream.Response().LogID())
		defer stream.Close()

		event, err := stream.Recv()
		as.Nil(err)
		as.Equal(WorkflowEventTypeInterrupt, event.Event)
		as.Equal("event1", event.Interrupt.InterruptData.EventID)
		as.Equal(1, event.Interrupt.InterruptData.Type)
		as.Equal("Question", event.Interrupt.NodeTitle)
	})
}

func TestWorkflowEventParsing(t *testing.T) {
	as := assert.New(t)
	t.Run("parse workflow event error", func(t *testing.T) {
		data := `{"error_code":400,"error_message":"Bad Request"}`
		event, err := ParseWorkflowEventError(data)
		as.Nil(err)
		as.Equal(400, event.ErrorCode)
		as.Equal("Bad Request", event.ErrorMessage)
	})

	t.Run("parse workflow event interrupt", func(t *testing.T) {
		data := `{"interrupt_data":{"event_id":"event1","type":1},"node_title":"Question"}`
		interrupt, err := ParseWorkflowEventInterrupt(data)
		as.Nil(err)
		as.Equal("event1", interrupt.InterruptData.EventID)
		as.Equal(1, interrupt.InterruptData.Type)
		as.Equal("Question", interrupt.NodeTitle)
	})

	t.Run("invalid json parsing", func(t *testing.T) {
		_, err := ParseWorkflowEventError("invalid json")
		as.Error(err)

		_, err = ParseWorkflowEventInterrupt("invalid json")
		as.Error(err)
	})
}
