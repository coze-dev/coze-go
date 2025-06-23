package coze

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowsChat(t *testing.T) {
	as := assert.New(t)
	t.Run("stream chat success", func(t *testing.T) {
		chat := newWorkflowsChat(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/workflows/chat", req.URL.Path)
			return mockStreamResponse(`event: conversation.chat.created
data: {"id":"chat1","conversation_id":"test_conversation_id","bot_id":"bot1","status":"created"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id","role":"assistant","content":"Hello"}

event: done
data: {}

`)
		})))
		stream, err := chat.Stream(context.Background(), &WorkflowsChatStreamReq{
			WorkflowID: "test_workflow",
			AdditionalMessages: []*Message{
				{
					Role:    MessageRoleUser,
					Content: "Hello",
				},
			},
			Parameters: map[string]any{
				"test": "value",
			},
		})
		as.NoError(err)
		as.NotNil(stream)
		as.NotEmpty(stream.Response().LogID())
		defer stream.Close()

		event1, err := stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationChatCreated, event1.Event)

		event2, err := stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationMessageDelta, event2.Event)

		event3, err := stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventDone, event3.Event)

		_, err = stream.Recv()
		as.Equal(io.EOF, err)
	})

	t.Run("Stream chat with error response", func(t *testing.T) {
		chat := newWorkflowsChat(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, baseResponse{
				Code: 100,
				Msg:  "Invalid workflow ID",
			})
		})))
		_, err := chat.Stream(context.Background(), &WorkflowsChatStreamReq{
			WorkflowID: "invalid_workflow",
		})
		fmt.Println(err)
		as.NotNil(err)

		cozeErr, ok := AsCozeError(err)
		as.True(ok)
		as.Equal(100, cozeErr.Code)
		as.Equal("Invalid workflow ID", cozeErr.Message)
		as.Equal("test_log_id", cozeErr.LogID)
	})
}
